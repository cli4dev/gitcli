package db

import (
	"fmt"
	"path/filepath"
)

//GetSQL 获取sql语句
func (tb *Table) GetSQL(outPath string) (fpath string, c string, err error) {
	fpath = filepath.Join(outPath, fmt.Sprintf("%s.sql", tb.Name))
	c, err = tb.translate(mysqlTmpl, tb.GetTmpt(""))
	return fpath, c, err
}

//GetGoFile 获取sql语句
func (tb *Table) GetGoFile(outPath string) (fpath string, c string, err error) {

	fpath = filepath.Join(outPath, fmt.Sprintf("%s.sql.go", tb.Name))
	c, err = tb.translate(mysqlTmpl, tb.GetTmpt(outPath))
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
