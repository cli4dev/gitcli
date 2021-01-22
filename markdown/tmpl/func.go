package tmpl

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
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

		//枚举处理函数
		"fIsEnumTB": hasKW("di", "dn"), //数据表的字段是否包含字典数据配置
		"fHasDT":    hasKW("dt"),       //数据表是否包含字典类型字段
		"fIsDI":     getKWS("di"),      //字段是否为字典ID
		"fIsDN":     getKWS("dn"),      //字段是否为字典Name
		"fIsDT":     getKWS("dt"),      //字段是否为字典Type

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
		"order":     getOrderBy,

		"ismysql":  stringsEqual("mysql"),
		"isoracle": stringsEqual("oracle"),
		"SL":       getKWS("sl"), //表单下拉框
		"CB":       getKWS("cb"), //表单复选框
		"RB":       getKWS("rb"), //表单单选框
		"TA":       getKWS("ta"), //表单文本域
		"DT":       getKWS("dt"), //表单日期选择器

		"query":     getRows("q"),                        //查询字段
		"list":      getRows("l"),                        //列表展示字段
		"detail":    getRows("r"),                        //详情展示字段
		"create":    getRows("c"),                        //创建字段
		"delete":    getRows("d"),                        //删除时判定字段
		"update":    getRows("u"),                        //更新字段
		"moduleCon": getBracketContent("sl", "cb", "rb"), //获取组件约束的内容
		"firstStr":  getStringByIndex(0),                 //获取约束的内容

		"rpath": getRouterPath,
		"fpath": getFilePath,

		"var":    getVar,
		"vars":   joinVars,
		"isTime": isTime,

		"lowerName": fGetLowerCase, //小驼峰式命名
		"upperName": fGetUpperCase, //大驼峰式命名
		//	"contains": contains,     //是否包含子串
		"lname": fGetLastName, //取最后一个单词
		"dpath": GetDetailPath,
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

func fGetLowerCase(n string) string {
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

func fGetUpperCase(n string) string {
	_, f := filepath.Split(n)
	f = strings.ReplaceAll(f, ".", "_")
	items := strings.Split(f, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}

func fGetLastName(n string) string {
	sp := "/"
	if strings.Contains(n, "_") {
		sp = "_"
	}

	names := strings.Split(strings.Trim(n, sp), sp)

	if len(names) > 2 {
		return strings.Join(names[2:len(names)], sp)
	}
	return names[len(names)-1]
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
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

func getOrderBy(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Rows))
	fileds := []string{}
	orders := []string{}
	ob := map[string]string{}

	for _, v := range tb.Rows {
		if strings.Contains(v.Con, "OB") {
			if !strings.Contains(v.Con, "OB(") {
				fileds = append(fileds, v.Name)
				continue
			}
			for _, v1 := range strings.Split(v.Con, ",") {
				if !strings.Contains(v1, "OB(") {
					continue
				}
				s := strings.Index(v1, "OB(")
				e := strings.Index(v1, ")")
				orders = append(orders, v1[s+1:e])
				ob[v1[s+1:e]] = v.Name
			}
		}
	}

	if len(orders) > 0 {
		sort.Sort(sort.StringSlice(orders))
	}

	for _, v := range orders {
		fileds = append(fileds, ob[v])
	}

	for _, v := range fileds {
		row := map[string]interface{}{
			"name":  v,
			"comma": true,
		}
		columns = append(columns, row)
	}

	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getSeqs() func(tb *Table) []map[string]interface{} {
	return func(tb *Table) []map[string]interface{} {
		columns := make([]map[string]interface{}, 0, len(tb.Rows))

		for _, v := range tb.Rows {
			if strings.Contains(v.Con, "SEQ") {
				descsimple := strings.Join(getBracketContent("SEQ")(v.Desc), ",")
				row := map[string]interface{}{
					"name":       v.Name,
					"descsimple": descsimple,
					"seqname":    "seqname", //  fmt.Sprintf("seq_%s_%s", fGetNName(tb.Name), getFilterName(tb.Name, v.Cname)),
					"desc":       v.Desc,
					"type":       v.Type,
					"len":        v.Len,
					"comma":      true,
				}
				columns = append(columns, row)
			}
		}
		return columns
	}
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

func stringsEqual(s string) func(s1 string) bool {
	return func(s1 string) bool {
		return strings.EqualFold(s, s1)
	}
}

//getRouterPath .
func getRouterPath(tabName string) string {
	if tabName == "" {
		return ""
	}
	return "/" + strings.Replace(strings.ToLower(tabName), "_", "/", -1)
}

//getFilePath .
func getFilePath(tabName string) string {
	if tabName == "" {
		return ""
	}
	return "/" + strings.Replace(strings.ToLower(tabName), "_", ".", -1)
}

//GetDetailPath .
func GetDetailPath(tabName string) string {
	dir, f := filepath.Split(strings.Replace(tabName, "_", "/", -1))
	return "/" + dir + f
}

func getStringByIndex(index int) func(s []string) string {
	return func(s []string) string {
		return types.GetStringByIndex(s, index)
	}
}

func getBracketContent(keys ...string) func(con string) []string {
	return func(con string) []string {
		s := ""
		for _, key := range keys {
			rex := regexp.MustCompile(fmt.Sprintf(`%s\((.+?)\)`, key))
			strs := rex.FindAllString(con, -1)
			if len(strs) < 1 {
				continue
			}
			str := strs[0]
			str = strings.TrimPrefix(str, fmt.Sprintf("%s(", key))
			str = strings.TrimRight(str, ")")
			s = fmt.Sprintf("%s,%s", s, str)
		}
		if s == "" {
			return []string{}
		}
		s = strings.TrimLeft(s, ",")
		return strings.Split(s, ",")
	}
}
func hasKW(tp ...string) func(t *Table) bool {
	return func(t *Table) bool {
		ext := map[string]bool{}
		for _, r := range t.Rows {
			for _, t := range tp {
				if isCons(r.Con, t) {
					ext[t] = true
				}
			}
		}
		for _, t := range tp {
			if _, ok := ext[t]; !ok {
				return false
			}
		}
		return true
	}
}
