package tmpts

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//
const (
	IN = "input"       //输入框
	SL = "select"      //下拉框
	CB = "checkbox"    //复选框
	RB = "radio"       //单选框
	TA = "textarea"    //文本域
	DT = "date-picker" //日期选择器
)

//GetInputData 获取模板数据
func getInputData(tb *Table) map[string]interface{} {
	return map[string]interface{}{
		"name":          tb.Name,
		"desc":          tb.Desc,
		"dblink":        tb.DBLink,
		"createcolumns": getCreateColumns(tb), //创建数据必需的字段
		"getcolumns":    getSingleDetail(tb),  //单条数据要显示的字段
		"querycolumns":  getQueryColumns(tb),  //查询字段
		"updatecolumns": getUpdateColumns(tb), //可更新的字段
		"selectcolumns": getListColumns(tb),   //列表要显示的字段
		"deletecolumns": getDeleteColumns(tb), //根据主键删除
		"pk":            getPks(tb),
		"seqs":          getSeqs(tb),
		"path":          GetRouterPath(tb.Name),
		"dpath":         GetDetailPath(tb.Name),
		"di":            getDictionariesID(tb),
		"dn":            getDictionariesName(tb),
	}
}

//GetInputData 获取模板数据
func getMDFileInputData(tb *Table) map[string]interface{} {
	return map[string]interface{}{
		"name":          getName(tb.Name),
		"tname":         tb.Name,
		"desc":          tb.Desc,
		"dblink":        tb.DBLink,
		"createcolumns": getCreateColumns(tb), //创建数据必需的字段
		"getcolumns":    getSingleDetail(tb),  //单条数据要显示的字段
		"querycolumns":  getQueryColumns(tb),  //查询字段
		"updatecolumns": getUpdateColumns(tb), //可更新的字段
		"selectcolumns": getListColumns(tb),   //列表要显示的字段
		"deletecolumns": getDeleteColumns(tb), //根据主键删除
		"pk":            getPks(tb),
		"seqs":          getSeqs(tb),
		"order":         getOrderBy(tb),
		"path":          GetRouterPath(tb.Name),
		"di":            getDictionariesID(tb),
		"dn":            getDictionariesName(tb),
	}
}

//getMDInputData 获取模板数据
func getMDInputData(tb *Table) map[string]interface{} {
	input := map[string]interface{}{
		"name":    tb.Name,
		"desc":    tb.Desc,
		"columns": getColumns(tb),
	}

	return input
}

func getName(n string) string {
	return strings.Replace(n, "_", "/", 2)
}

func getDictionariesID(tb *Table) string {
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "DI") {
			return v.Cname
		}
	}
	return ""
}

func getDictionariesName(tb *Table) string {
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "DN") {
			return v.Cname
		}
	}
	return ""
}

//GetRouterPath .
func GetRouterPath(tabName string) string {
	return "/" + strings.Replace(tabName, "_", "/", -1)
}

//GetDetailPath .
func GetDetailPath(tabName string) string {
	dir, f := filepath.Split(strings.Replace(tabName, "_", "/", -1))
	return "/" + dir + f
}

func getSingleDetail(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if tagExist(v.Con, "R") {
			descsimple := removeBracket(v.Desc)
			name := v.Cname
			row := map[string]interface{}{
				"name":       name,
				"pname":      v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
				"domType":    getDomType(v.Con, v.Type),
				"source":     getSource(v.Cname, v.Con),
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getDomType(cons, t string) string {
	if strings.Contains(t, "date") || strings.Contains(cons, "DT") {
		return DT
	}
	if strings.Contains(cons, "SL") {
		return SL
	}
	if strings.Contains(cons, "CB") {
		return CB
	}
	if strings.Contains(cons, "RB") {
		return RB
	}
	if strings.Contains(cons, "TA") {
		return TA
	}

	return IN
}

func isFormBox(v string) bool {
	return strings.Contains(v, "SL") || strings.Contains(v, "CB") || strings.Contains(v, "RB")
}

//前端页面获取相应的字典数据
func getSource(name, cons string) map[string]interface{} {
	m := make(map[string]interface{})
	var con string
	var is bool
	for _, v := range strings.Split(cons, ",") {
		if isFormBox(v) {
			con = v
			is = true
			break
		}
	}
	if !is {
		return nil
	}
	if !strings.Contains(con, "(") || !strings.Contains(con, ")") {
		m[name] = map[string]string{
			"path":   "",
			"params": "",
		}
		return m
	}
	s := strings.Index(con, "(")
	e := strings.Index(con, ")")
	p := strings.Split(con[s+1:e], ",")
	path := GetRouterPath(p[0]) + "/getdictionary"
	c := map[string]string{
		"path":   path,
		"params": "{}",
	}
	m[name] = c
	return m
}

func getColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		row := map[string]interface{}{
			"name":   v.Cname,
			"def":    v.Def,
			"type":   v.Type,
			"isnull": v.IsNull,
			"cons":   v.Con,
			"len":    v.Len,
			"lenstr": v.LenStr,
			"desc":   v.Desc,
			"comma":  true,
		}
		columns = append(columns, row)

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getCreateColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if tagExist(v.Con, "C") && !strings.Contains(v.Con, "SEQ") {
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
				"domType":    getDomType(v.Con, v.Type),
				"source":     getSource(v.Cname, v.Con),
				"isnull":     v.IsNull,
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getQueryColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if tagExist(v.Con, "Q") && !strings.Contains(v.Con, "SEQ") {
			descsimple := removeBracket(v.Desc)

			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
				"domType":    getDomType(v.Con, v.Type),
				"source":     getSource(v.Cname, v.Con),
			}
			columns = append(columns, row)
		}

	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getUpdateColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if tagExist(v.Con, "U") && !strings.Contains(v.Con, "SEQ") && !strings.Contains(v.Con, "PK") {
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
				"domType":    getDomType(v.Con, v.Type),
				"source":     getSource(v.Cname, v.Con),
				"isnull":     v.IsNull,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getDeleteColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		//仅根据pk删除
		if tagExist(v.Con, "D") && strings.Contains(v.Con, "PK") {
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

//获取列表字段
func getListColumns(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if tagExist(v.Con, "L") {
			name := v.Cname
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       name,
				"pname":      v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"domType":    getDomType(v.Con, v.Type),
				"comma":      true,
				"source":     getSource(v.Cname, v.Con),
				"isquery":    tagExist(v.Con, "Q"),
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getPks(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "PK") {
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getSeqs(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))

	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "SEQ") {
			descsimple := removeBracket(v.Desc)
			row := map[string]interface{}{
				"name":       v.Cname,
				"descsimple": descsimple,
				"seqname":    fmt.Sprintf("seq_%s_%s", fGetNName(tb.Name), getFilterName(tb.Name, v.Cname)),
				"desc":       v.Desc,
				"type":       v.Type,
				"len":        v.Len,
				"comma":      true,
			}
			columns = append(columns, row)
		}
	}
	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getOrderBy(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Columns))
	fileds := []string{}
	orders := []string{}
	ob := map[string]string{}

	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "OB") {
			if !strings.Contains(v.Con, "OB(") {
				fileds = append(fileds, v.Cname)
				continue
			}
			for _, v1 := range strings.Split(v.Con, ",") {
				if !strings.Contains(v1, "OB(") {
					continue
				}
				s := strings.Index(v1, "OB(")
				e := strings.Index(v1, ")")
				orders = append(orders, v1[s+1:e])
				ob[v1[s+1:e]] = v.Cname
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

func fGetNName(n string) string {
	items := strings.Split(n, "_")
	if len(items) <= 1 {
		return n
	}
	return strings.Join(items[1:], "_")
}

func getFilterName(t string, f string) string {
	text := make([]string, 0, 1)
	tb := strings.Split(t, "_")
	fs := strings.Split(f, "_")
	for _, v := range fs {
		ex := false
		for _, k := range tb {
			if v == k {
				ex = true
				break
			}
		}
		if !ex {
			text = append(text, v)
		}
	}
	if len(text) == 0 {
		return "id"
	}
	return strings.Join(text, "_")

}

func gGetModulePackageName(module []string) []string {
	npkgs := make([]string, 0, len(module)/2)
	n := make(map[string]string)
	for _, m := range module {
		nm := fGetServicePackagePath(m)
		if _, ok := n[nm]; !ok {
			npkgs = append(npkgs, nm)
			n[nm] = nm
		}
	}
	return npkgs
}

func tagExist(con, tag string) bool {
	s := strings.Split(con, ",")
	for _, v := range s {
		if isPropersubset("CDLQRU", v) && strings.Contains(v, tag) {
			return true
		}
	}
	return false
}

//判断b是不是a的真子集
func isPropersubset(a string, b string) bool {
	var hash int
	var c rune = 'A'
	start := uint(c)
	for k := range []rune(a) {
		hash |= (1 << (uint(a[k]) - start))
	}

	for k := range []rune(b) {
		if hash&(1<<(uint(b[k])-start)) == 0 {
			return false
		}
	}
	return true
}

func removeBracket(s string) string {
	if strings.Contains(s, "(") {
		s = s[:strings.Index(s, "(")]
	}
	if strings.Contains(s, "（") {
		s = s[:strings.Index(s, "（")]
	}
	return s
}
