package gitlabs

import (
	"encoding/json"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/lib4go/envs"
	"github.com/micro-plat/lib4go/sysinfo/pipes"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//Repository 仓库信息
type Repository struct {
	Name     string `json:"name"`
	Desc     string `json:"description"`
	Path     string `json:"relative_path"`
	Type     string `json:"type"`
	FullPath string `json:"-"`
}

//NewRepository 创建仓库
func NewRepository(fullPath string) *Repository {
	u, _ := url.Parse(fullPath)
	return &Repository{FullPath: fullPath, Path: u.Path, Type: "project"}
}

//String 输出内容
func (r *Repository) String() string {
	if buff, err := json.Marshal(&r); err == nil {
		return string(buff)
	}
	return ""
}

//GetLocalPath 获取本地路径
func (r *Repository) GetLocalPath() string {
	u, _ := url.Parse(r.FullPath)
	gopath := envs.GetString("GOPATH")
	return filepath.Join(gopath, "src", u.Host, r.Path)
}
func (r *Repository) GetShortPath(s string) string {
	gopath := envs.GetString("GOPATH")
	src := filepath.Join(gopath, "src")
	return strings.Replace(r.GetLocalPath(), src, "~", -1)
}

//Exists 本地仓库是否是存在
func (r *Repository) Exists() bool {
	rpath := filepath.Join(r.GetLocalPath(), ".git")
	if _, err := os.Stat(rpath); err != nil {
		return os.IsExist(err)
	}
	return true
}

//Reset 拉取项目
func (r *Repository) Reset(branch ...string) error {
	session := sh.InteractiveSession()
	session.SetDir(r.GetLocalPath())
	for _, b := range branch {
		session.Command("git", "branch")
		buff, err := session.Output()
		if err != nil {
			return err
		}
		if hasBranch(string(buff), b) {

			logs.Log.Infof("%s > git reset --hard", r.GetLocalPath())
			session.Command("git", "reset", "--hard")
			if err := session.Run(); err != nil {
				return err
			}

			logs.Log.Infof("%s > git checkout %s", r.GetLocalPath(), b)
			session.Command("git", "checkout", b)
			if err := session.Run(); err != nil {
				return err
			}
			logs.Log.Infof("%s > git reset --hard", r.GetLocalPath())
			session.Command("git", "reset", "--hard")
			if err := session.Run(); err != nil {
				return err
			}
		}

	}
	return nil
}

//Pull 拉取项目
func (r *Repository) Pull(branch ...string) error {
	session := sh.InteractiveSession()
	session.SetDir(r.GetLocalPath())
	session.Command("git", "branch")
	buff, err := session.Output()
	if err != nil {
		return err
	}
	for _, b := range branch {
		if hasBranch(string(buff), b) {
			logs.Log.Infof("%s > git pull origin %s", r.GetLocalPath(), b)
			session.Command("git", "pull", "origin", b)
		} else {
			logs.Log.Infof("%s > git pull origin %s:%s", r.GetLocalPath(), b, b)
			session.Command("git", "pull", "origin", b+":"+b)
		}
		if err := session.Run(); err != nil {
			return err
		}
	}
	return nil
}

//Clone 克隆项目
func (r *Repository) Clone() error {
	cmd := fmt.Sprintf(`git clone %s %s`, r.FullPath, r.GetLocalPath())
	_, err := pipes.RunString(cmd)
	if err != nil && strings.Contains(err.Error(), "exit status 128") {
		return fmt.Errorf("fatal: 目标路径 '%s' 已经存在，并且不是一个空目录。", r.GetLocalPath())
	}
	return err
}
func hasBranch(s string, b string) bool {
	items := strings.Split(s, "\n")
	for _, i := range items {
		if strings.Contains(i, b) {
			return true
		}
	}
	return false
}
