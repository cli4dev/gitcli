package tmpl

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

const MYSQL = "mysql"

type callHanlder func(string) string

func getfuncs(tp string) map[string]interface{} {
	return map[string]interface{}{
		"pascal":   fPascal,      //获取pascal变量名称
		"rmhd":     rmhd,         //去除首段名称
		"isnull":   isNull(tp),   //返回空语句
		"short":    shortWord,    //获取特殊字段前的字符串
		"sql":      sqlType(tp),  //转换为SQL的数据类型
		"codeType": codeType,     //转换为GO代码的数据类型
		"def":      defValue(tp), //返回SQL中的默认值
		"seq":      getSEQ(tp),   //获取SEQ的变量值
		"pks":      getPKS,       //获取主键列表
		"maxIndex": getMaxIndex,
		"indexs":   getDBIndex(tp), //获取表的索引串
		"lower":    getLower,       //获取变量的最小写字符
	}
}

func getLower(s string) string {
	return strings.ToLower(s)
}

//去掉首段名称
func rmhd(input string) string {
	index := strings.Index(input, "_")
	return input[index+1:]
}

//获取短文字
func shortWord(input string) string {
	reg := regexp.MustCompile(`^[\p{Han}|\w]+`)
	return reg.FindString(input)
}

//获取短文字
func isNull(tp string) callHanlder {
	switch tp {
	case MYSQL:
		return func(input string) string {
			return mysqlIsNull[input]
		}
	}
	return func(input string) string { return "" }
}

//首字母大写，并去掉下划线
func fPascal(input string) string {
	items := strings.Split(input, "_")
	nitems := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if len(item) == 1 {
			nitems = append(nitems, strings.ToUpper(item[0:1]))
			continue
		}
		if strings.EqualFold(item, "id") {
			nitems = append(nitems, strings.ToUpper(item))
			continue
		}
		nitems = append(nitems, strings.ToUpper(item[0:1])+item[1:])
	}
	return strings.Join(nitems, "")
}

//通过正则表达式，转换正确的数据库类型
func sqlType(tp string) callHanlder {
	switch tp {
	case MYSQL:
		return func(input string) string {
			buff := []byte(strings.ToLower(input))
			for k, v := range tp2mysql {
				reg := regexp.MustCompile(k)
				if reg.Match(buff) {
					if !strings.Contains(v, "*") {
						return v
					}
					value := reg.FindStringSubmatch(input)
					if len(value) > 1 {
						return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
					}
					return v
				}
			}
			return input
		}
	}
	return func(input string) string { return "" }
}

//通过正则表达式，转换正确的数据库类型
func codeType(input string) string {
	buff := []byte(strings.ToLower(input))
	for k, v := range any2code {
		reg := regexp.MustCompile(k)
		if reg.Match(buff) {
			return v
		}
	}
	return input
}

//通过正则表达式，转换正确的数据库类型
func defValue(tp string) callHanlder {
	switch tp {

	case MYSQL:
		return func(input string) string {
			buff := []byte(strings.Trim(strings.ToLower(input), "'"))
			for _, defs := range def2mysql {
				for k, v := range defs {
					reg := regexp.MustCompile(k)
					if reg.Match(buff) {
						if !strings.Contains(v, "*") {
							return v
						}
						value := reg.FindStringSubmatch(input)
						if len(value) > 1 {
							return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
						}
						return input
					}
				}
			}
			return input
		}
	}
	return func(input string) string { return "" }
}
func getPKS(t *Table) []string {
	return t.GetPKS()
}
func getSEQ(tp string) func(r *Row) string {
	switch tp {
	case MYSQL:
		return func(r *Row) string {
			if isCons(r.Con, "seq") {
				return "AUTO_INCREMENT"
			}
			return ""
		}
	}
	return func(r *Row) string { return "" }
}
func getMaxIndex(r interface{}) int {
	v := reflect.ValueOf(r)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array || v.Kind() == reflect.Map {
		return v.Len() - 1
	}
	return 0
}

//去掉首段名称
func isCons(input string, tp string) bool {
	cks, ok := cons[tp]
	if !ok {
		return false
	}
	buff := []byte(strings.ToLower(input))
	for _, ck := range cks {
		reg := regexp.MustCompile(ck)
		if reg.Match(buff) {
			return true
		}
	}
	return false
}

func getPKINdex(r *Table) string {
	keys := make([]string, 0, 1)
	for _, r := range r.Rows {
		if isCons(r.Con, "pk") {
			keys = append(keys, r.Name)
		}
	}
	if len(keys) == 0 {
		return ""
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(keys, ","))

}
func getDBIndex(tp string) func(r *Table) string {
	switch tp {
	case MYSQL:
		return func(r *Table) string {
			indexs := r.GetIndexs()
			list := make([]string, 0, len(indexs))
			for _, index := range indexs {
				list = append(list, fmt.Sprintf("KEY %s(%s)", index.Name, index.fields.Join(",")))
			}
			pks := getPKINdex(r)
			if pks != "" {
				list = append(list, pks)
			}
			if len(list) > 0 {
				return "," + strings.Join(list, ",")
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

func getIndex(input string) (bool, string, int) {
	buff := []byte(strings.Trim(strings.ToLower(input), "'"))
	for _, v := range cons["idx"] {
		reg := regexp.MustCompile(v)
		if reg.Match(buff) {
			value := reg.FindStringSubmatch(strings.ToLower(input))
			if len(value) > 2 {
				return true, value[1], types.GetInt(value[2], 0)
			}
			if len(value) > 1 {
				return true, value[1], 0
			}
		}
	}
	return false, "", 0
}
