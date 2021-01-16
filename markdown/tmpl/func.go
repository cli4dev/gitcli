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
		"varName": getVarName, //获取pascal变量名称
		"names":   getNames,
		"mod":     getMod,
		"rmhd":    rmhd,       //去除首段名称
		"isNull":  isNull(tp), //返回空语句

		"isEnumTB": isEnumTB,
		"isDI":     getKWS("di"),
		"isDN":     getKWS("dn"),

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

		"SL": getKWS("sl"), //表单下拉框
		"CB": getKWS("cb"), //表单复选框
		"RB": getKWS("rb"), //表单单选框
		"TA": getKWS("ta"), //表单文本域
		"DT": getKWS("dt"), //表单日期选择器

		"query":  getRows("q"),            //查询字段
		"list":   getRows("l"),            //列表展示字段
		"detail": getRows("r"),            //详情展示字段
		"create": getRows("c"),            //创建字段
		"delete": getRows("d"),            //删除时判定字段
		"update": getRows("u"),            //更新字段
		"SLCon":  getBracketContent("sl"), //获取约束的内容

		"rpath": getRouterPath,

		"var":    getVar,
		"vars":   joinVars,
		"isTime": isTime,
	}
}

func getLower(s string) string {
	return strings.ToLower(s)
}
func getMod(x int, y int) int {
	return x % y
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
		if strings.EqualFold(item, "id") || strings.EqualFold(item, "url") {
			nitems = append(nitems, strings.ToUpper(item))
			continue
		}
		nitems = append(nitems, strings.ToUpper(item[0:1])+item[1:])
	}
	return strings.Join(nitems, "")
}
func getNames(input string) []string {
	items := strings.Split(strings.Trim(input, "_"), "_")
	return items
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
		cks = cons["*"]
	}
	buff := []byte(strings.ToLower(input))
	for _, ck := range cks {
		nck := types.DecodeString(strings.Contains(ck, "%s"), true, fmt.Sprintf(ck, tp), ck)
		reg := regexp.MustCompile(nck)
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

func getRows(tp ...string) func(row []*Row) []*Row {
	return func(row []*Row) []*Row {
		list := make([]*Row, 0, 1)
		for _, r := range row {
		NEXT:
			for _, t := range tp {
				if isCons(r.Con, t) {
					list = append(list, r)
					break NEXT
				}
			}
		}
		return list
	}
}

func getKWS(tp ...string) func(input string) bool {
	return func(input string) bool {
		for _, t := range tp {
			if isCons(input, t) {
				return true
			}
		}
		return false
	}
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

func getJoin(text ...string) string {
	return strings.Join(text, "")
}

var vars = map[string][]string{}

func joinVars(name string) []string {
	return vars[name]
}

func getVar(name string, value ...string) string {
	if len(value) == 0 {
		return strings.Join(vars[name], "")
	}
	if t, ok := vars[name]; ok {
		old := make([]string, 0, len(t)+len(value))
		old = append(old, t...)
		old = append(old, value...)
		vars[name] = old
	} else {
		vars[name] = value
	}
	return ""
}
func isTime(input string) bool {
	tp := codeType(input)
	return tp == "time.Time"
}

//getRouterPath .
func getRouterPath(tabName string) string {
	if tabName == "" {
		return ""
	}
	return "/" + strings.Replace(strings.ToLower(tabName), "_", "/", -1)
}

func getBracketContent(key string) func(con string) []string {
	return func(con string) []string {
		rex := regexp.MustCompile(fmt.Sprintf(`%s\((.+?)\)`, key))
		strs := rex.FindAllString(con, -1)
		if len(strs) < 1 {
			return []string{""}
		}
		str := strs[0]
		str = strings.TrimPrefix(str, fmt.Sprintf("%s(", key))
		str = strings.TrimRight(str, ")")
		return strings.Split(str, ",")
	}
}
func isEnumTB(t *Table) bool {
	var di, dn = false, false
	for _, r := range t.Rows {
		if isCons(r.Con, "di") {
			di = true
		}
		if isCons(r.Con, "dn") {
			dn = true
		}
		if di && dn {
			return true
		}
	}
	return false
}
