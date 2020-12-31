package db

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/micro-plat/lib4go/types"
)

func getPkg(path string) string {
	if path == "" {
		return ""
	}
	names := strings.Split(strings.Trim(path, "/"), "/")
	pkgName := names[len(names)-1]
	return pkgName
}

//GetTmpt 获取转换模板
func (tb *Table) GetTmpt(outpath string) map[string]interface{} {
	pkg := getPkg(outpath)
	columns := make([]map[string]interface{}, 0, len(tb.Columns))
	for i, v := range tb.Columns {
		cType, err := tb.fGetType(v.Type, tb.DBType)
		if err != nil {
			return nil
		}

		descExt := types.DecodeString(v.DescExt, "", "", fmt.Sprintf("(%s)", v.DescExt))
		columns = append(columns, map[string]interface{}{
			"name":     v.Cname,
			"desc":     strings.Replace(v.Desc, ";", " ", -1),
			"desc_ext": descExt,
			"type":     cType,
			"len":      v.Len,
			"def":      tb.getDef(v.Def, v.Con, cType),
			"seq":      tb.getSeq(v.Con),
			"null":     tb.getNull(v.IsNull),
			"not_end":  i < len(tb.Columns)-1,
		})
	}
	uks, err := tb.getUniqueKeys()
	if err != nil {
		return nil
	}
	keys, err := tb.getKeys()
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"name":           tb.Name,
		"pkg":            pkg,
		"desc":           tb.Desc,
		"columns":        columns,
		"uks":            uks,
		"pk":             tb.getPk(),
		"auto_increment": tb.getAutoIncrement(),
		"keys":           keys,
	}
}

func (tb *Table) Translate(c string, input interface{}) (string, error) {
	var tmpl = template.New("mysql").Funcs(tb.makeFunc())
	np, err := tmpl.Parse(c)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return strings.Replace(buff.String(), "{###}", "`", -1), nil
}

func (tb *Table) getNull(n bool) string {
	if n {
		return ""
	}
	return "not null"
}

func (tb *Table) getDef(n string, c string, ctype string) string {
	if strings.Contains(c, "SEQ") {
		return ""
	}
	if n == "" || strings.ToLower(n) == "null" {
		return ""
	}
	if strings.TrimSpace(n) == "-" {
		return "default '-'"
	}
	if strings.TrimSpace(n) == "sysdate" {
		return "default current_timestamp"
	}

	if strings.Contains(strings.ToLower(ctype), "varchar") {
		return fmt.Sprintf("default '%s'", n)
	}

	return "default " + n
}

//getUniqueKeys 名最长64
func (tb *Table) getUniqueKeys() ([]map[string]string, error) {
	uks := make([]map[string]string, 0)
	orderUks := make(map[string][]string) //带顺序的uk
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "UNQ") {
			for _, con := range strings.SplitAfter(v.Con, "),") {
				params := getBracketContent(con, "UNQ")
				if len(params) > 2 {
					return nil, fmt.Errorf("%s的UNQ格式不正确,%+v", v.Cname, params)
				}
				if len(params) == 0 {
					ukName := fmt.Sprintf("unq_%s_%s", tb.Name, v.Cname)
					uks = append(uks, map[string]string{
						"ukname":  ukName,
						"ukfield": v.Cname,
					})
					continue
				}
				if len(params) == 1 {
					if _, ok := orderUks[params[0]]; ok {
						orderUks[params[0]] = append(orderUks[params[0]], v.Cname)
						continue
					}
					orderUks[params[0]] = []string{v.Cname}
					continue
				}
				order, err := strconv.Atoi(params[1])
				if err != nil {
					return nil, fmt.Errorf("%s的UNQ格式的顺序不正确:%+v", v.Cname, err)
				}
				field := fmt.Sprintf("%d:%s", order, v.Cname)
				if _, ok := orderUks[params[0]]; ok {
					orderUks[params[0]] = append(orderUks[params[0]], field)
					continue
				}
				orderUks[params[0]] = []string{field}
			}
		}
	}

	for k, v := range orderUks {
		sort.Sort(sort.StringSlice(v))
		fields := []string{}
		for _, val := range v {
			f := strings.Split(val, ":")
			fields = append(fields, f[len(f)-1])
		}
		uks = append(uks, map[string]string{
			"ukname":  k,
			"ukfield": strings.Join(fields, ","),
		})
	}
	return uks, nil
}

//getKeys 名最长64
func (tb *Table) getKeys() ([]map[string]string, error) {
	keys := []map[string]string{}
	orderKeys := map[string][]string{} //带顺序的KEY
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "IDX") {
			for _, con := range strings.SplitAfter(v.Con, "),") {
				params := getBracketContent(con, "IDX")
				if len(params) > 2 {
					return nil, fmt.Errorf("%s的IDX格式不正确,%+v", v.Cname, params)
				}
				if len(params) == 0 {
					ukName := fmt.Sprintf("key_%s_%s", tb.Name, v.Cname)
					keys = append(keys, map[string]string{
						"kname":  ukName,
						"kfield": v.Cname,
					})
					continue
				}
				if len(params) == 1 {
					if _, ok := orderKeys[params[0]]; ok {
						orderKeys[params[0]] = append(orderKeys[params[0]], v.Cname)
						continue
					}
					orderKeys[params[0]] = []string{v.Cname}
					continue
				}
				//if len(params) == 2
				order, err := strconv.Atoi(params[1])
				if err != nil {
					return nil, fmt.Errorf("%s的IDX格式的顺序不正确:%+v", v.Cname, err)
				}
				field := fmt.Sprintf("%d:%s", order, v.Cname)
				if _, ok := orderKeys[params[0]]; ok {
					orderKeys[params[0]] = append(orderKeys[params[0]], field)
					continue
				}
				orderKeys[params[0]] = []string{field}
			}
		}

	}

	for k, v := range orderKeys {
		sort.Sort(sort.StringSlice(v))
		fields := []string{}
		for _, val := range v {
			f := strings.Split(val, ":")
			fields = append(fields, f[len(f)-1])
		}
		keys = append(keys, map[string]string{
			"kname":  k,
			"kfield": strings.Join(fields, ","),
		})
	}
	return keys, nil
}
func (tb *Table) getPk() string {
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "PK") {
			return v.Cname
		}
	}
	return ""
}

func (tb *Table) getSeq(v string) string {
	if strings.Contains(v, "SEQ") {
		return "AUTO_INCREMENT"
	}
	return ""
}

func (tb *Table) getAutoIncrement() string {
	for _, v := range tb.Columns {
		if strings.Contains(v.Con, "SEQ") {
			strs := getBracketContent(v.Con, "SEQ")
			if len(strs) == 1 {
				return fmt.Sprintf("AUTO_INCREMENT=%s", strs[0])
			}
			return ""
		}
	}
	return ""
}

func (tb *Table) makeFunc() map[string]interface{} {
	// nfuncs := funcs
	// nfuncs[]
	//  map[string]interface{}{return
	// 	"cName": tb.fGetCName,
	// 	"nName": tb.fGetNName,
	// 	"sub1":  tb.sub1,
	// }
	return funcs
}
func (tb *Table) fGetCName(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}
func (tb *Table) fGetNName(n string) string {
	items := strings.Split(n, "_")
	if len(items) <= 1 {
		return n
	}
	return strings.Join(items[1:], "_")
}

func (tb *Table) fGetType(n string, dbType string) (string, error) {
	if dbType == "mysql" {
		return n, nil
	}

	return ConvertDataType(n)
}

func (tb *Table) getFilterName(t string, f string) string {
	text := make([]string, 0, 1)
	ti := strings.Split(t, "_")
	fs := strings.Split(f, "_")
	for _, v := range fs {
		ex := false
		for _, k := range ti {
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

func getBracketContent(s string, key string) []string {
	rex := regexp.MustCompile(fmt.Sprintf(`%s\((.+?)\)`, key))
	strs := rex.FindAllString(s, -1)
	if len(strs) < 1 {
		return nil
	}
	str := strs[0]
	str = strings.TrimPrefix(str, fmt.Sprintf("%s(", key))
	str = strings.TrimRight(str, ")")
	return strings.Split(str, ",")
}
