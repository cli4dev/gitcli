package db

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type Table struct {
	Name       string //表名
	Desc       string //表描述
	DBLink     string //dblink
	Class      string //分类
	Columns    TableColumn
	Addition   bool   //是否是附加
	DBType     string //数据类型
	SystemName string //系统名字
	MaxSort    int    //最大排序值
	Menu       string //前端菜单名
}

type Column struct {
	Cname   string //字段名
	Len     int    //长度
	LenStr  string //长度
	Type    string //类型
	Def     string //默认值
	IsNull  bool   //为空
	Con     string //约束
	Desc    string //描述
	DescExt string //描述的括号里面的内容
	Sort    int    //排序
	After   string
}

func (t *Table) String() string {
	buff := strings.Builder{}
	buff.WriteString(t.Name)
	buff.WriteString("(")
	buff.WriteString(t.Desc)
	buff.WriteString(")")
	buff.WriteString("\n")
	for _, c := range t.Columns {
		buff.WriteString(fmt.Sprintf("%s\t%s\n", c.Cname, c.Type))

	}
	return buff.String()
}

//TableColumn 表的列排序用
type TableColumn []*Column

func (t TableColumn) Len() int {
	return len(t)
}

//从低到高
func (t TableColumn) Less(i, j int) bool {
	if t[i].Sort < t[j].Sort {
		return true
	}
	if t[i].Sort == t[j].Sort && t[i].After == "" {
		return true
	}
	return false
}

func (t TableColumn) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func NewTable(name, desc, dblink string) *Table {
	return &Table{
		Name:    name,
		Desc:    desc,
		DBLink:  dblink,
		Columns: TableColumn{},
	}
}

func (t *Table) SetClass(c string) {
	t.Class = c
}

func (t *Table) AppendColumn(c *Column) error {
	if strings.Contains(c.Desc, "(") && strings.Contains(c.Desc, ")") {
		desc := c.Desc
		c.Desc = desc[:strings.Index(desc, "(")]
		c.DescExt = desc[strings.Index(desc, "(")+1 : strings.Index(desc, ")")]
	}
	t.Columns = append(t.Columns, c)
	return nil
}

func FilterTable(tables []*Table, filters []string) []*Table {
	newTables := make([]*Table, len(tables))
	if len(filters) > 0 {
		for _, tb := range tables {
			e := false
			for _, f := range filters {
				if strings.EqualFold(tb.Name, f) {
					e = true
					break
				}
			}
			if !e {
				newTables = append(tables, tb)
			}
		}
		return newTables
	}
	newTables = tables
	return newTables
}

func (t *Table) Mysql2Column(data types.XMap) (c *Column) {
	return &Column{
		Cname:  data.GetString("column_name"),
		Len:    -1,
		LenStr: "",
		Type:   strings.ToLower(data.GetString("column_type")),
		Def:    getDefault(data.GetString("column_default")),
		IsNull: getBool(data.GetString("is_nullable")),
		Con:    getMysqlCon(data.GetString("con")),
		Desc:   data.GetString("column_comment"),
	}
}

func (t *Table) Oracle2Column(data types.XMap) (c *Column) {
	con := fmt.Sprintf("%s,%s", data.GetString("constraint_type"), data.GetString("index_name"))
	return &Column{
		Cname:  strings.ToLower(data.GetString("column_name")),
		Len:    data.GetInt("data_length", -1),
		LenStr: data.GetString("data_length"),
		Type:   strings.ToLower(data.GetString("data_type")),
		Def:    getDefault(data.GetString("data_default")),
		IsNull: getBool(data.GetString("nullable")),
		Con:    getOracleCon(con),
		Desc:   strings.TrimSpace(strings.ReplaceAll(data.GetString("column_comments"), "|", "&#124;")),
	}
}

func getMysqlCon(con string) string {
	if strings.Contains(con, "PRIMARY KEY") {
		con = "PK"
	}
	if strings.Contains(con, "UNIQUE") {
		con = strings.ReplaceAll(con, "UNIQUE", "UNQ")
	}
	return con
}

// C　　　　　　Check constraint on a table
// P　　　　　　Primary key
// U　　　　　　Unique key
// R　　　　　　Referential integrity
// V　　　　　　With check option, on a view
// O　　　　　　With read only, on a view
// H　　　　　　Hash expression
// F　　　　　　Constraint that involves a REF column
// S　　　　　　Supplemental loggin
func getOracleCon(s string) string {
	cons := strings.Split(s, ",")
	con := []string{}
	for _, v := range cons {
		if strings.HasPrefix(v, "C(") { //检查约束，跳过
			continue
		}
		if strings.HasPrefix(v, "P(") {
			con = append(con, "PK")
			continue
		}
		if strings.HasPrefix(v, "U(") {
			con = append(con, strings.Replace(strings.Replace(v, "U(", "UNQ(", 1), "|", ",", 1))
			continue
		}
		if strings.HasPrefix(v, "IDX(") {
			con = append(con, strings.Replace(v, "|", ",", 1))
		}
	}
	return strings.Join(con, ",")
}

func getDefault(str string) string {
	s := strings.TrimSpace(strings.ReplaceAll(str, "|", "&#124;"))
	s = strings.TrimPrefix(s, `'`)
	s = strings.TrimSuffix(s, `'`)
	return s
}

func getBool(in interface{}) bool {
	switch in.(string) {
	case "NO", "no", "N", "n":
		return false
	case "YES", "yes", "Y", "y":
		return true
	}
	return false
}
