package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

//getHomePath 获取用户home目录 没实现跨平台
func getHomePath() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.HomeDir
}

//GetProjectPath 获取项目路径
func GetProjectPath(projectPath string) (string, string, error) {
	npath := projectPath
	if !strings.HasPrefix(npath, "./") && !strings.HasPrefix(npath, "/") && !strings.HasPrefix(npath, "../") {
		srcPath, err := os.Getwd()
		if err != nil {
			return "", "", err
		}
		npath = filepath.Join(srcPath, npath)
	}

	root, err := filepath.Abs(npath)
	if err != nil {
		return "", "", fmt.Errorf("不是有效的项目路径:%s", root)
	}
	return filepath.Base(root), root, nil
}

//GetWebSrcPath 获取web项目src目录
//判断路径下是否有src目录且src下有App.vue,有则返回src目录和项目目录
//默认返回空
func GetWebSrcPath(projectPath string) (string, string) {
	n := strings.LastIndex(projectPath, "src")
	if n < 0 {
		return "", ""
	}
	parentDir := projectPath[0:n]
	srcPath := path.Join(parentDir, "src")
	appVuePath := path.Join(srcPath, "App.vue")
	if pathExists(appVuePath) { //存在返回
		return parentDir, srcPath
	}
	return GetWebSrcPath(parentDir)
}

//GetProjectBasePath 如果开启了gomod 则返回module名
//未使用gomod则判断path中是否存在$GOPATH，存在则返回$GOPATH下面的名字
//默认返回空
func GetProjectBasePath(projectPath string) (string, error) {
	cmd := exec.Command("go", "env")
	envs := []byte{}
	envs, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行go env出错，%+v", err)
	}
	var basePath, gopath, gomod string
	for _, v := range strings.Split(string(envs), "\n") {
		if strings.HasPrefix(v, "GOPATH=") {
			gopath = strings.TrimPrefix(v, `GOPATH="`)
			gopath = strings.TrimRight(gopath, `"`)
			continue
		}
		if strings.HasPrefix(v, "GOMOD=") {
			gomod = strings.TrimPrefix(v, `GOMOD="`)
			gomod = strings.TrimRight(gomod, `"`)
		}
	}

	if gomod != "" && strings.Contains(gomod, projectPath) {
		f, err := os.Open(gomod)
		if err != nil {
			return "", fmt.Errorf("打开%s文件出错，%+v", gomod, err)
		}
		defer f.Close()

		br := bufio.NewReader(f)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			line := string(a)
			if strings.HasPrefix(line, "module ") {
				basePath = strings.TrimPrefix(line, "module ")
				break
			}
		}
		return basePath, nil
	}

	if gopath != "" {
		if strings.Contains(projectPath, gopath) {
			basePath = strings.TrimPrefix(projectPath, fmt.Sprintf("%s/src/", gopath))
		}
		return basePath, nil
	}
	return "", nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
