package tmpts

import (
	"fmt"
	"path/filepath"
)

//GetScheme 获取数据库结构
func (tb *Table) GetScheme(outPath string, toGoFile bool) (fpath string, c string, err error) {
	if !toGoFile {
		fpath = filepath.Join(outPath, fmt.Sprintf("%s.sql", tb.Name))
		c, err = tb.Translate(mysqlTmpl, tb.GetTmpt(""))
		return fpath, c, err
	}
	fpath = filepath.Join(outPath, fmt.Sprintf("%s.sql.go", tb.Name))
	c, err = tb.Translate(mysqlTmpl, tb.GetTmpt(outPath))
	return fpath, c, err
}

//GetDBInstallFile 获取db install文件
func (tb *MarkDownDB) GetDBInstallFile(outPath string) (path string, c string, err error) {
	pkgName := getPkg(outPath)
	path = filepath.Join(outPath, "install.go")
	input := map[string]interface{}{
		"pkg": pkgName,
		"tbs": tb.Tables,
	}
	c, err = tb.translate(mysqlInstallTmpl, input)
	return path, c, err
}

//GetEntity 获取实体
func (tb *Table) GetEntity() (c string, err error) {
	return tb.Translate(entityTmpl, tb.GetTmpt(""))
}

//GetSelect 获取select查询语句
func (tb *Table) GetSelect() (c string, err error) {
	c, err = tb.Translate(selectSingle, tb.GetTmpt(""))
	return c, err
}
