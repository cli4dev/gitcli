package tmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//Table 表名称
type Table struct {
	Name   string //表名
	Desc   string //表描述
	PKG    string //包名称
	Rows   []*Row
	Indexs Indexs
}

//Row 行信息
type Row struct {
	Name   string //字段名
	Type   string //类型
	Def    string //默认值
	IsNull string //为空
	Con    string //约束
	Desc   string //描述
}

//Indexs 索引集
type Indexs map[string]*Index

//Index 索引
type Index struct {
	fields fields
	Name   string
}
type fields []*Field

//Field 字段信息
type Field struct {
	Name  string
	Index int
}

func (t fields) Len() int {
	return len(t)
}
func (t fields) Join(s string) string {
	list := make([]string, 0, len(t))
	for _, v := range t {
		list = append(list, v.Name)
	}
	return strings.Join(list, s)
}

//从低到高
func (t fields) Less(i, j int) bool {
	if t[i].Index < t[j].Index {
		return true
	}
	if t[i].Index == t[j].Index {
		return true
	}
	return false
}

func (t fields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

//NewTable 创建表
func NewTable(name, desc string) *Table {
	return &Table{
		Name: name,
		Desc: desc,
		Rows: make([]*Row, 0, 1),
	}
}

//AddRow 添加行信息
func (t *Table) AddRow(r *Row) error {
	t.Rows = append(t.Rows, r)
	return nil
}

//SetPkg 添加行信息
func (t *Table) SetPkg(path string) {
	names := strings.Split(strings.Trim(path, "/"), "/")
	t.PKG = names[len(names)-1]
}

//GetPKS 获取主键列表
func (t *Table) GetPKS() []string {
	list := make([]string, 0, 1)
	for _, r := range t.Rows {
		if isCons(r.Con, "pk") {
			list = append(list, r.Name)
		}
	}
	return list
}

//GetIndexs 获取所有索引信息
func (t *Table) GetIndexs() Indexs {
	if t.Indexs != nil {
		return t.Indexs
	}
	indexs := map[string]*Index{}
	for ri, r := range t.Rows {
		ok, name, index := getIndex(r.Con)
		if !ok {
			continue
		}
		index = types.DecodeInt(index, 0, ri)
		if v, ok := indexs[name]; ok {
			v.fields = append(v.fields, &Field{Name: r.Name, Index: index})
			continue
		}
		indexs[name] = &Index{Name: name, fields: []*Field{{Name: r.Name, Index: index}}}
	}
	for _, index := range indexs {
		sort.Sort(index.fields)
	}
	t.Indexs = indexs
	return t.Indexs
}
func (t *Table) String() string {
	buff := strings.Builder{}
	buff.WriteString(t.Name)
	buff.WriteString("(")
	buff.WriteString(t.Desc)
	buff.WriteString(")")
	buff.WriteString("\n")
	for _, c := range t.Rows {
		buff.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%v\t%s\n", c.Name, c.Type, c.Con, c.Def, c.IsNull, c.Desc))

	}
	return buff.String()
}

//Translate 翻译模板
func Translate(c string, tp string, input interface{}) (string, error) {
	var tmpl = template.New("table").Funcs(getfuncs(tp))
	np, err := tmpl.Parse(c)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return strings.Replace(strings.Replace(buff.String(), "{###}", "`", -1), "&#39;", "'", -1), nil
}
