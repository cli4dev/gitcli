package tmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

//Table 表名称
type Table struct {
	Name    string //表名
	Desc    string //表描述
	Columns TableColumn
}

//TableColumn 表的列排序用
type TableColumn []*Column

//Column 列信息
type Column struct {
	Cname  string //字段名
	Len    string //长度
	ILen   int    //长度
	Type   string //类型
	Def    string //默认值
	IsNull bool   //为空
	Con    string //约束
	Desc   string //描述
}

//NewTable 创建表
func NewTable(name, desc string) *Table {
	return &Table{
		Name:    name,
		Desc:    desc,
		Columns: TableColumn{},
	}
}

//AddColumn 添加列
func (t *Table) AddColumn(c *Column) error {
	t.Columns = append(t.Columns, c)
	return nil
}
func (t *Table) String() string {
	buff := strings.Builder{}
	buff.WriteString(t.Name)
	buff.WriteString("(")
	buff.WriteString(t.Desc)
	buff.WriteString(")")
	buff.WriteString("\n")
	for _, c := range t.Columns {
		buff.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%v\t%s\n", c.Cname, c.Type, c.Con, c.Def, c.IsNull, c.Desc))

	}
	return buff.String()
}

//Translate 翻译模板
func (t *Table) Translate(c string, input map[string]interface{}) (string, error) {
	var tmpl = template.New("table").Funcs(funcs)
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
