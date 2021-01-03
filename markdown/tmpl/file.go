package tmpl

import (
	"fmt"
	"os"
	"path/filepath"
)

//GetSchemePath 获取Scheme路径
func GetSchemePath(outpath string, name string, gofile bool) string {
	path := filepath.Join(outpath, fmt.Sprintf("%s.sql", name))
	if gofile {
		path = filepath.Join(outpath, fmt.Sprintf("%s.sql.go", name))

	}
	return path
}

//GetInstallPath 获取DB安装文件
func GetInstallPath(outpath string) string {
	return filepath.Join(outpath, "install.go")
}

//Create 创建文件，文件夹 存在时写入则覆盖
func Create(path string, append bool) (file *os.File, err error) {
	dir := filepath.Dir(path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			err = fmt.Errorf("创建文件夹%s失败:%v", path, err)
			return nil, err
		}
	}

	_, err = os.Stat(path)
	var srcf *os.File
	if os.IsNotExist(err) {
		srcf, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
			return nil, err
		}
		return srcf, nil

	}
	if !append {
		return nil, fmt.Errorf("文件:%s已经存在", path)
	}
	srcf, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
		return nil, err
	}
	return srcf, nil

}
