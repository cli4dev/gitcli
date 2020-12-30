package db

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type mysqlTable struct {
}

//GetSQL 获取sql语句
func GetSQL(tbs []*Table, outPath string, pkg string) (tmpls map[string]string, err error) {
	t := &mysqlTable{}
	tmpls = make(map[string]string, len(tbs))
	for _, tb := range tbs {
		columns := make([]map[string]interface{}, 0, len(tb.Columns))
		for i, v := range tb.Columns {
			cType, err := t.fGetType(v.Type, tb.DBType)
			if err != nil {
				return nil, err
			}
			descExt := ""
			if v.DescExt != "" {
				descExt = fmt.Sprintf("(%s)", v.DescExt)
			}
			columns = append(columns, map[string]interface{}{
				"name":     v.Cname,
				"desc":     strings.Replace(v.Desc, ";", " ", -1),
				"desc_ext": descExt,
				"type":     cType,
				"len":      v.Len,
				"def":      t.getDef(v.Def, v.Con, cType),
				"seq":      t.getSeq(v.Con),
				"null":     t.getNull(v.IsNull),
				"not_end":  i < len(tb.Columns)-1,
			})
		}
		uks, err := t.getUniqueKeys(tb)
		if err != nil {
			return nil, err
		}
		keys, err := t.getKeys(tb)
		if err != nil {
			return nil, err
		}
		fpath := filepath.Join(outPath, fmt.Sprintf("%s.sql", tb.Name))
		if pkg != "" {
			fpath = filepath.Join(outPath, fmt.Sprintf("%s.sql.go", tb.Name))
		}
		tmpls[fpath], err = t.translate(DBMysqlTmpl, map[string]interface{}{
			"name":           tb.Name,
			"pkg":            pkg,
			"desc":           tb.Desc,
			"columns":        columns,
			"uks":            uks,
			"pk":             t.getPk(tb),
			"auto_increment": t.getAutoIncrement(tb),
			"keys":           keys,
		})
		if err != nil {
			return nil, err
		}
	}
	return tmpls, nil
}

func (t *mysqlTable) translate(c string, input interface{}) (string, error) {
	var tmpl = template.New("mysql").Funcs(t.makeFunc())
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

func (t *mysqlTable) getNull(n bool) string {
	if n {
		return ""
	}
	return "not null"
}

func (t *mysqlTable) getDef(n string, c string, ctype string) string {
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

//名最长64
func (t *mysqlTable) getUniqueKeys(tb *Table) ([]map[string]string, error) {
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
				//if len(params) == 2
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

//名最长64
func (t *mysqlTable) getKeys(tb *Table) ([]map[string]string, error) {
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
func (t *mysqlTable) getPk(tbs *Table) string {
	for _, v := range tbs.Columns {
		if strings.Contains(v.Con, "PK") {
			return v.Cname
		}
	}
	return ""
}

func (t *mysqlTable) getSeq(v string) string {
	if strings.Contains(v, "SEQ") {
		return "AUTO_INCREMENT"
	}
	return ""
}

func (t *mysqlTable) getAutoIncrement(tbs *Table) string {
	for _, v := range tbs.Columns {
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

func (t *mysqlTable) makeFunc() map[string]interface{} {
	return map[string]interface{}{
		"cName": t.fGetCName,
		"nName": t.fGetNName,
		"sub1":  t.sub1,
	}
}
func (t *mysqlTable) fGetCName(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}
func (t *mysqlTable) fGetNName(n string) string {
	items := strings.Split(n, "_")
	if len(items) <= 1 {
		return n
	}
	return strings.Join(items[1:], "_")
}

func (t *mysqlTable) fGetType(n string, dbType string) (string, error) {
	if dbType == "mysql" {
		return n, nil
	}

	return ConvertDataType(n)
}

func (m *mysqlTable) getFilterName(t string, f string) string {
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

func (t *mysqlTable) sub1(n int) int {
	return n - 1
}
