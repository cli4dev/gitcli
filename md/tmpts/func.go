package tmpts

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

var funcs = makeFunc()

func makeFunc() map[string]interface{} {
	return map[string]interface{}{
		"aname":      fGetAName,
		"dName":      fGetDName,
		"cname":      fGetCName,
		"pname":      fGetPName,
		"ctype":      fGetType,
		"cstype":     fsGetType,
		"lname":      fGetLastName,
		"lower":      fToLower,
		"vname":      vName,
		"puname":     pathUpperName,          //路径大写
		"humpName":   fGetHumpName,           //多个单词首字符大写
		"spkgName":   fGetServicePackageName, //包路径
		"mpkgName":   fGetModulePackageName,  //包路径
		"lName":      fGetLastName,           //取最后一个单词
		"fName":      fGetFirstName,          //取第一个单词
		"fServer":    fServer,                //判断是否有这个服务
		"getAppconf": getAppconf,
		"checkName":  checkName,      //是否含有sql
		"fUpperName": fGetFUpperName, //首字母大写
		"mName":      gGetMNAME,      //service找到对应的module
		"sqlName":    fGetSQLName,    //sql名
		"mod":        getMod,         //取余数
		"sub1":       sub1,
		"rmheader":   fRMHeader,
	}
}

func vName(v string) string {
	items := strings.Split(v, "/")
	nitems := make([]string, 0, len(items))
	for k, i := range items {
		if k == 0 {
			nitems = append(nitems, i)
		}
		if k > 0 {
			nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
		}

	}
	return strings.Join(nitems, "")
}

func pathUpperName(v string) string {
	items := strings.Split(v, "/")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}

func fGetAName(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for k, i := range items {
		if k == 0 {
			nitems = append(nitems, i)
		}
		if k > 0 {
			nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
		}

	}
	return strings.Join(nitems, "")
}
func fGetDName(n string) string {
	if n == "" {
		n = "default"
	}
	return n
}
func fGetCName(n string) string {
	_, f := filepath.Split(n)
	f = strings.ReplaceAll(f, ".", "_")
	items := strings.Split(f, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		if i == "id" {
			nitems = append(nitems, strings.ToUpper(i))
			continue
		}
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}

func fGetPName(n string) string {
	dir, _ := filepath.Split(n)
	items := strings.Split(dir, "/")
	if len(items) >= 2 {
		return items[len(items)-2]
	}
	return "modules"
}

func fGetSQLName(n string) string {
	n = strings.ReplaceAll(n, "/", "_")
	return fGetCName(n)
}

func fRMHeader(n string) string {
	items := strings.Split(n, "_")
	return strings.Join(items[1:], "_")
}

func fGetType(v string) string {
	n := strings.ToLower(v)
	switch {
	case strings.Contains(n, "varchar"):
		return "string"
	case strings.Contains(n, "number"):
		var i, j int
		for k, v := range n {
			if v == '(' {
				i = k
			}
			if v == ')' {
				j = k
			}
		}
		ii, _ := strconv.Atoi(n[i+1 : j])
		if strings.Contains(n, ",") {
			j = strings.Index(n, ",")
			if ii <= 10 {
				return "float32"
			}
			return "float64"
		}
		if ii <= 10 {
			return "int"
		}
		return "int64"

	case strings.Contains(n, "date"):
		return "time.Time"
	default:
		fmt.Println("default:", n)
		return "string"
	}
}
func fsGetType(n string) string {
	return fGetType(n)
}

func fGetFUpperName(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}
func fGetLastName(n string) string {
	sp := "/"
	if strings.Contains(n, "_") {
		sp = "_"
	}

	names := strings.Split(strings.Trim(n, sp), sp)

	if len(names) > 2 {
		return strings.Join(names[2:], sp)
	}
	return names[len(names)-1]
}

func fToLower(s string) string {
	return strings.ToLower(s)
}

func fServer(s, substr string) bool {
	return strings.Contains(s, substr)
}

func fGetFirstName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	return names[0]
}

func fGetHumpName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	buff := bytes.NewBufferString("")
	for _, v := range names {
		buff.WriteString(fGetLoopHumpName(v, "."))
	}
	return strings.Replace(buff.String(), ".", "", -1)
}

func fGetLoopHumpName(n string, s string) string {
	names := strings.Split(strings.Trim(n, s), s)
	buff := bytes.NewBufferString("")
	for _, v := range names {
		buff.WriteString(strings.ToUpper(v[0:1]))
		buff.WriteString(v[1:])
	}
	return strings.Replace(buff.String(), ".", "", -1)
}

func fGetServicePackageName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	if len(names) == 1 {
		return "services"
	}
	return strings.ToLower(names[len(names)-2])
}

func fGetPackageName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[0:len(names)-1], "/")
}

func fGetModulePackageName(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	if len(names) == 1 {
		return "modules"
	}
	return strings.ToLower(names[len(names)-2])
}

func fGetServicePackagePath(n string) string {
	names := strings.Split(strings.Trim(n, "/"), "/")
	if len(names) == 1 {
		return "services"
	}
	return strings.ToLower(filepath.Join("services", strings.Join(names[0:len(names)-1], "/")))
}

func getAppconf(str string, index int) string {
	strArray := strings.Split(str, "|")
	if len(strArray) < index {
		return ""
	}
	if ok := strArray[index-1]; ok != "" {
		return ok
	}
	return ""
}

func checkName(str string) bool {
	return strings.Contains(str, "sql")
}

func gGetMNAME(str string) string {
	dir, _ := filepath.Split(str)
	if strings.EqualFold(dir, "") {
		return ""
	}
	return fmt.Sprintf("/%s", strings.TrimRight(dir, "/"))
}

func getMod(a, b int) int {
	return a % b
}
func sub1(n int) int {
	return n - 1
}
