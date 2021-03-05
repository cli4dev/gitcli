package server

import (
	"fmt"
	ifs "io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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

	s.fs.Start()
	go s.start()

	for {
		select {
		case <-s.notifyChan:
			s.pause()
			go s.start()
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
	s.running = true

	//开启文件监控
	go s.watch()

	//文件打包
	err := s.session.Command("go", "install").Run()
	if err != nil {
		return
	}

	//程序启动
	s.session.Command(s.serverName, "run").Run()

	return
}

func (s *server) pause() {
	if s.running {
		s.running = false
		s.session.Kill(os.Interrupt)
	}
}

func (s *server) close() (err error) {

	var sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	select {
	case <-sigChan:
	case err = <-s.errChan:
		s.running = false
		close(s.closeChan)
		close(s.notifyChan)
		close(s.errChan)
		time.Sleep(time.Second)
	}

	return err
}

var defExcludePath = []string{"vendor", "node_modules", ".gitignore", ".gitcli"}

func (s *server) isExclude(path string) bool {
	for _, v := range defExcludePath {
		if strings.Contains(path, v) {
			return true
		}
	}
	return false
}

func (s *server) watch() {
	filepath.WalkDir(s.path, func(path string, d ifs.DirEntry, err error) error {
		if d.IsDir() && !s.isExclude(path) {
			go s.watchChildren(path)
		}
		return nil
	})
}

func (s *server) watchChildren(path string) {
	//监控子节点变化
	ch, err := s.fs.WatchChildren(path)
	if err != nil {
		s.fs.Close()
		s.errChan <- err
		return
	}

	for {
		select {
		case <-s.ticker.C:
			if s.hasNotify {
				s.notifyChan <- 1
				s.hasNotify = false
				return
			}
		case <-s.closeChan:
			s.fs.Close()
			return
		case cldWatcher := <-ch:
			if cldWatcher.GetError() != nil && s.running {
				s.errChan <- fmt.Errorf("监控项目文件发生错误：%+v", cldWatcher.GetError())
				return
			}
			logs.Log.Info("----------------------项目发生变化，应用程序重启----------------------")
			if !s.isExclude(cldWatcher.GetPath()) {
				s.hasNotify = true
			}
		LOOP:
			ch, err = s.fs.WatchChildren(path)
			if err != nil {
				if s.running {
					s.errChan <- fmt.Errorf("应用程序运行中，未获取到文件监控")
					return
				}
				logs.Log.Errorf("文件监控错误%+v", err)
				time.Sleep(time.Second * 3)
				goto LOOP
			}
		}
	}
}
