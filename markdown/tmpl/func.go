package tmpl

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//MYSQL mysql数据库
const MYSQL = "mysql"

type callHanlder func(string) string

func getfuncs(tp string) map[string]interface{} {
	return map[string]interface{}{
		"varName":   getVarName,      //获取pascal变量名称
		"rmhd":      rmhd,            //去除首段名称
		"isNull":    isNull(tp),      //返回空语句
		"shortName": shortName,       //获取特殊字段前的字符串
		"dbType":    dbType(tp),      //转换为SQL的数据类型
		"codeType":  codeType,        //转换为GO代码的数据类型
		"defValue":  defValue(tp),    //返回SQL中的默认值
		"seqTag":    getSEQTag(tp),   //获取SEQ的变量值
		"seqValue":  getSEQValue(tp), //获取SEQ起始值
		"pks":       getPKS,          //获取主键列表
		"indexs":    getDBIndex(tp),  //获取表的索引串
		"maxIndex":  getMaxIndex,     //最大索引值
		"lower":     getLower,        //获取变量的最小写字符
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
func shortName(input string) string {
	reg := regexp.MustCompile(`^[\p{Han}|\w]+`)
	return reg.FindString(input)
}

//获取短文字
func isNull(tp string) func(*Row) string {
	switch tp {
	case MYSQL:
		return func(row *Row) string {
			return mysqlIsNull[row.IsNull]
		}
	}
	return func(row *Row) string { return "" }
}

//首字母大写，并去掉下划线
func getVarName(input string) string {
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
func dbType(tp string) callHanlder {
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
func defValue(tp string) func(*Row) string {
	switch tp {

	case MYSQL:
		return func(row *Row) string {
			if isCons(row.Con, "seq") {
				return ""
			}
			buff := []byte(strings.Trim(strings.ToLower(row.Def), "'"))
			for _, defs := range def2mysql {
				for k, v := range defs {
					reg := regexp.MustCompile(k)
					if reg.Match(buff) {
						if !strings.Contains(v, "*") {
							return v
						}
						value := reg.FindStringSubmatch(row.Def)
						if len(value) > 1 {
							return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
						}
						return row.Def
					}
				}
			}
			return row.Def
		}
	}
	return func(row *Row) string { return "" }
}
func getPKS(t *Table) []string {
	return t.GetPKS()
}
func getSEQTag(tp string) func(r *Row) string {
	switch tp {
	case MYSQL:
		return func(r *Row) string {
			if isCons(r.Con, "seq") {
				return "auto_increment"
			}
			return ""
		}
	}
	return func(r *Row) string { return "" }
}
func getSEQValue(tp string) func(r *Table) string {
	switch tp {
	case MYSQL:
		return func(r *Table) string {
			for _, r := range r.RawRows {
				if isCons(r.Con, "seq") {
					if v := types.GetInt(r.Def, 0); v != 0 {
						return fmt.Sprintf("auto_increment = %d", v)
					}

				}
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
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
	cks, ok := cons[strings.ToLower(tp)]
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

func getDBIndex(tp string) func(r *Table) string {
	switch tp {
	case MYSQL:
		return func(r *Table) string {
			indexs := r.GetIndexs()
			list := make([]string, 0, len(indexs))
			for _, index := range indexs {
				switch index.Type {
				case "idx":
					list = append(list, fmt.Sprintf("index %s(%s)", index.Name, index.fields.Join(",")))
				case "unq":
					list = append(list, fmt.Sprintf("unique index %s(%s)", index.Name, index.fields.Join(",")))
				case "pk":
					list = append(list, fmt.Sprintf("primary key (%s)", index.fields.Join(",")))
				}
			}
			if len(list) > 0 {
				return "," + strings.Join(list, "\n\t\t,")
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

func getIndex(input string, tp string) (bool, string, int) {
	buff := []byte(strings.Trim(strings.ToLower(input), "'"))
	for _, v := range cons[tp] {
		reg := regexp.MustCompile(v)
		if reg.Match(buff) {
			value := reg.FindStringSubmatch(strings.ToLower(input))
			if len(value) > 2 {
				return true, value[1], types.GetInt(value[2], 0)
			}
			if len(value) > 1 {
				return true, value[1], 0
			}
			return true, "", 0
		}
	}
	return false, "", 0
}

//getKWCons 获取关键字约束列
func getKWCons(input string, keyword string) bool {
	for _, kw := range keywordMatch {
		reg := regexp.MustCompile(fmt.Sprintf(kw, keyword))
		if reg.Match([]byte(input)) {
			return true
		}
	}
	return false
}
