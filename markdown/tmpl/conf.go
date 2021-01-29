package tmpl

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/micro-plat/lib4go/security/md5"
)

//SnippetConf 用于vue的路由,hydra服务的注册,import的路径等代码片段生成
type SnippetConf struct {
	Name      string `json:"name"`
	HasDetail bool   `json:"has_detail"`
	BasePath  string `json:"base_path"`
	Desc      string `json:"desc"`
}

func NewSnippetConf(t *Table) *SnippetConf {
	rows := getRows("r")(t.Rows)
	return &SnippetConf{
		Name:      t.Name,
		HasDetail: len(rows) > 0,
		BasePath:  t.BasePath,
		Desc:      t.Desc,
	}
}

func (t *SnippetConf) SaveConf(confPath string) error {
	if confPath == "" {
		return nil
	}

	s, err := Read(confPath)
	if err != nil {
		return err
	}

	//创建文件
	fs, err := Create(confPath, true)
	if err != nil {
		return err
	}

	//设置
	conf := make(map[string]interface{}, 0)
	if len(s) > 0 {
		err = json.Unmarshal(s, &conf)
		if err != nil {
			return err
		}
	}
	conf[t.Name] = t
	r, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	fs.WriteString(string(r))
	fs.Close()
	return nil
}

func GetSnippetConf(path string) ([]*SnippetConf, error) {

	s, err := Read(path)
	if err != nil {
		return nil, err
	}

	conf := make(map[string]*SnippetConf, 0)
	if len(s) > 0 {
		err = json.Unmarshal(s, &conf)
		if err != nil {
			return nil, err
		}
	}
	confs := make([]*SnippetConf, 0)
	for _, v := range conf {
		confs = append(confs, v)
	}

	return confs, nil
}

//FieldConf 用于field文件生成
type FieldConf struct {
	Fields []FieldItem `json:"fields"`
}

type FieldItem struct {
	Desc  string `json:"desc"`
	Name  string `json:"name"`
	Table string `json:"table"`
}

func NewFieldConf(t *Table) *FieldConf {
	fields := []FieldItem{}
	for _, v := range t.Rows {
		item := FieldItem{
			Desc:  v.Desc,
			Name:  v.Name,
			Table: t.Name,
		}
		fields = append(fields, item)
	}
	return &FieldConf{Fields: fields}
}

func GetFieldConf(path string) (map[string]*FieldItem, error) {

	s, err := Read(path)
	if err != nil {
		return nil, err
	}

	conf := make(map[string]*FieldItem, 0)
	if len(s) > 0 {
		err = json.Unmarshal(s, &conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func (t *FieldConf) SaveConf(confPath string) error {
	if confPath == "" {
		return nil
	}

	s, err := Read(confPath)
	if err != nil {
		return err
	}

	//创建文件
	fs, err := Create(confPath, true)
	if err != nil {
		return err
	}

	//设置
	conf := make(map[string]interface{}, 0)
	if len(s) > 0 {
		err = json.Unmarshal(s, &conf)
		if err != nil {
			return err
		}
	}
	for _, v := range t.Fields {
		conf[v.Name] = v
	}

	r, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	fs.WriteString(string(r))
	fs.Close()
	return nil
}

func GetFieldConfPath(root string) string {
	projectName, projectPath, err := utils.GetProjectPath(root)
	if err != nil {
		panic(err)
	}
	if projectPath == "" {
		return ""
	}
	return path.Join(utils.GetGitcliHomePath(), fmt.Sprintf("server/%s_filed_%s.json", projectName, md5.Encrypt(projectPath)))
}

func GetVueConfPath(root string) string {
	_, projectPath, err := utils.GetProjectPath(root)
	if err != nil {
		panic(err)
	}
	webPath, _ := utils.GetWebSrcPath(projectPath)
	if webPath == "" {
		return ""
	}
	return path.Join(utils.GetGitcliHomePath(), fmt.Sprintf("web/web_%s.json", md5.Encrypt(webPath)))
}

func GetGoConfPath(root string) string {
	projectName, projectPath, err := utils.GetProjectPath(root)
	if err != nil {
		panic(err)
	}
	if projectPath == "" {
		return ""
	}
	return path.Join(utils.GetGitcliHomePath(), fmt.Sprintf("server/%s_%s.json", projectName, md5.Encrypt(projectPath)))
}
