package server

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/urfave/cli"
)

func runServer() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		//判断项目是否存在
		projectPath := utils.GetProjectPath(c.Args().Get(0))
		if !utils.PathExists(filepath.Join(projectPath, "main.go")) {
			return fmt.Errorf("未指定的运行应用程序的项目路径")
		}

		//构建服务
		s, err := newServer(filepath.Base(projectPath), projectPath)
		if err != nil {
			return err
		}

		//服务启动
		errChan := make(chan error, 1)
		go func() {
			if err := s.resume(); err != nil {
				errChan <- err
			}
		}()

		//文件监控
		go func() {
			if err := s.watch(); err != nil {
				errChan <- err
			}
		}()

		//服务退出
		var sigChan = make(chan os.Signal, 3)
		signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)
		select {
		case <-sigChan:
			s.close()
		case err = <-errChan:
			s.close()
			return err
		}

		return nil
	}
}
