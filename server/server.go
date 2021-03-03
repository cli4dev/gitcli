package server

import (
	"fmt"
	"os"
	"time"

	"github.com/codeskyblue/go-sh"
	logs "github.com/lib4dev/cli/logger"
)

type server struct {
	session    *sh.Session
	serverName string
	path       string
	fs         *fs
	running    bool
	notifyChan chan int
	closeChan  chan int
}

func newServer(serverName, path string) (*server, error) {
	session := sh.InteractiveSession()
	session.SetDir(path)
	r, err := NewFileSystem(path)
	if err != nil {
		return nil, err
	}
	return &server{
		serverName: serverName,
		path:       path,
		fs:         r,
		session:    session,
		notifyChan: make(chan int, 1),
		closeChan:  make(chan int, 1),
	}, nil
}

//Reset 拉取项目
func (s *server) resume() error {

	go s.start()

	for {
		select {
		case <-s.notifyChan:
			go func() {
				if s.running {
					s.pause()
				}
				s.start()
			}()
		case <-s.closeChan:
			s.pause()
			return nil
		}
	}
}

func (s *server) start() {
	if s.running {
		return
	}
	logs.Log.Info("进行应用程序安装")
	if err := s.session.Command("go", "install").Run(); err != nil {
		logs.Log.Error(err)
		return
	}
	logs.Log.Info("应用程序启动")
	s.running = true
	if err := s.session.Command(s.serverName, "run").Run(); err != nil {
		logs.Log.Error(err)
		return
	}
}

func (s *server) pause() {
	logs.Log.Info("关闭正在运行的应用程序")
	s.running = false
	s.session.Kill(os.Interrupt)
}

func (s *server) close() {
	s.running = false
	close(s.closeChan)
	close(s.notifyChan)
	time.Sleep(time.Second)
}

func (s *server) watch() error {

	s.fs.Start()

	//监控子节点变化
	ch, err := s.fs.WatchChildren(s.path)
	if err != nil {
		s.fs.Close()
		return err
	}

	deadline := time.Minute
	for {
		select {
		case <-time.After(deadline):
			if !s.running {
				return fmt.Errorf("超时未获取到文件监控")
			}
			//logs.Log.Info("项目未发生变化")
		case <-s.closeChan:
			s.fs.Close()
			logs.Log.Info("关闭文件监控")
			return nil
		case cldWatcher := <-ch:
			if cldWatcher.GetError() != nil {
				return fmt.Errorf("监控项目文件发生错误：%+v", cldWatcher.GetError())
			}
			logs.Log.Info("项目发生变化，应用程序重启")
			s.notifyChan <- 1
		LOOP:
			ch, err = s.fs.WatchChildren(s.path)
			if err != nil {
				if !s.running {
					return fmt.Errorf("应用程序未启动，未获取到文件监控")
				}
				logs.Log.Errorf("文件监控错误%+v", err)
				time.Sleep(time.Second * 5)
				goto LOOP
			}
		}
	}

}
