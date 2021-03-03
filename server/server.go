package server

import (
	"fmt"
	ifs "io/fs"
	"os"
	"path/filepath"
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
	hasNotify  bool
	notifyChan chan int
	closeChan  chan int
	ticker     *time.Ticker
	errChan    chan error
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
		errChan:    make(chan error, 1),
		ticker:     time.NewTicker(time.Second),
	}, nil
}

//Reset 拉取项目
func (s *server) resume() {

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
			return
		}
	}
}

func (s *server) start() {
	if s.running {
		return
	}
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
	close(s.errChan)
	time.Sleep(time.Second)
}

func (s *server) watch() {

	s.fs.Start()

	filepath.WalkDir(s.path, func(path string, d ifs.DirEntry, err error) error {
		if d.IsDir() { //@todo 排除不监控的文件
			go func() {
				if err := s.watchChildren(path); err != nil {
					s.errChan <- err
				}
			}()
		}
		return nil
	})

}

func (s *server) watchChildren(path string) error {
	//监控子节点变化
	ch, err := s.fs.WatchChildren(path)
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
		case <-s.ticker.C:
			if s.hasNotify {
				logs.Log.Info("项目发生变化，应用程序重启")
				s.notifyChan <- 1
			}
			s.hasNotify = false
		case <-s.closeChan:
			s.fs.Close()
			return nil
		case cldWatcher := <-ch:
			if cldWatcher.GetError() != nil {
				return fmt.Errorf("监控项目文件发生错误：%+v", cldWatcher.GetError())
			}
			fmt.Println("-----------", path)
			s.hasNotify = true
		LOOP:
			ch, err = s.fs.WatchChildren(path)
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
