package db

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type MarkDownDB struct {
	Tables []*Table
}

type Line struct {
	Text   string
	LineID int
}

type TableLine struct {
	Lines [][]*Line
}

func (tb *MarkDownDB) translate(c string, input interface{}) (string, error) {
	var tmpl = template.New("mysql")
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

//TableAddition 对表进行附加处理
func (m *MarkDownDB) TableAddition() error {
	baseTables := make(map[string]*Table, len(m.Tables))
	addTables := make([]*Table, 0, len(m.Tables))
	for _, v := range m.Tables {
		if v.Addition {
			addTables = append(addTables, v)
			continue
		}
		if _, ok := baseTables[v.Name]; ok {
			return fmt.Errorf("存在同名表：%s", v.Name)
		}
		baseTables[v.Name] = v
	}

	if len(addTables) > 0 {
		for _, v := range addTables {
			if _, ok := baseTables[v.Name]; !ok {
				return fmt.Errorf("附加表%s不存在相应的基础表", v.Name)
			}
			bt := baseTables[v.Name]
			if bt.DBType != v.DBType {
				return fmt.Errorf("附加表%s(%s)与基础表(%s)的数据库类型不正确", v.Name, v.DBType, bt.DBType)
			}
			//附加
			//基础表列
			columns := make(map[string]*Column, len(bt.Columns))
			for _, v1 := range bt.Columns {
				columns[v1.Cname] = v1
			}
			//替换
			maxSort := bt.MaxSort
			for _, v1 := range v.Columns {
				//找排序值
				v1.Sort += maxSort
				columns[v1.Cname] = v1
			}
			baseTables[v.Name].Columns = make([]*Column, 0, len(columns))
			for _, v1 := range columns {
				baseTables[v.Name].Columns = append(baseTables[v.Name].Columns, v1)
			}
		}

		for k, v := range baseTables {
			sorts := make(map[string]int, len(v.Columns))
			for _, v1 := range v.Columns {
				sorts[v1.Cname] = v1.Sort
			}
			for k1, v1 := range v.Columns {
				if v1.After != "" {
					if _, ok := sorts[v1.After]; ok {
						baseTables[k].Columns[k1].Sort = sorts[v1.After]
					}
				}
			}
		}

		m.Tables = make([]*Table, 0, len(baseTables))
		for _, v := range baseTables {
			m.Tables = append(m.Tables, v)
		}
	}
	return nil
}

//Markdown2DB 读取markdown文件并转换为MarkDownDB对象
func Markdown2DB(fn string) (*MarkDownDB, error) {
	lines, err := readMarkdown(fn)
	if err != nil {
		return nil, err
	}
	a := strings.Split(fn, ".")
	DbType := a[len(a)-2]
	SystemName := getSystemName(lines)

	tables, err := tableLine2Table(line2TableLine(lines), getTableClasses(lines), DbType, SystemName)
	if err != nil {
		return nil, err
	}

	db := &MarkDownDB{
		Tables: tables,
	}

	return db, nil
}

//readMarkdown 读取md文件
func readMarkdown(name string) ([]*Line, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := make([]*Line, 0, 64)
	rd := bufio.NewReader(f)
	num := 0
	for {
		num++
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Trim(line, "\n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, &Line{Text: line, LineID: num})
	}
	return lines, nil
}

func getSystemName(lines []*Line) string {
	for _, line := range lines {
		text := strings.TrimSpace(strings.Replace(line.Text, " ", "", -1))
		if strings.HasPrefix(text, "#") && strings.Count(text, "#") == 1 {
			return strings.TrimPrefix(text, "#")
		}
	}
	return ""
}

func getTableClasses(lines []*Line) map[string]string {
	var start bool
	classes := make(map[string]string)
	for _, line := range lines {
		text := strings.TrimSpace(strings.Replace(line.Text, " ", "", -1))
		if text == "|前缀|类型|" {
			start = true
			continue
		}
		if start && strings.Count(text, "|") != 3 {
			break
		}
		if start {
			foo := strings.Split(text, "|")
			classes[foo[1]] = foo[2]
		}
	}
	return classes
}

//lines2TableLine 数据行转变为以表为单个整体的数据行
func line2TableLine(lines []*Line) (tl TableLine) {
	//获取划分行divideLines
	dlines := []int{}
	for i, line := range lines {
		text := strings.TrimSpace(strings.Replace(line.Text, " ", "", -1))
		if text == "|字段名|类型|默认值|为空|约束|描述|" {
			dlines = append(dlines, i-1)
		}
		if len(dlines)%2 == 1 && strings.Count(text, "|") != 7 {
			dlines = append(dlines, i-1)
		}
	}
	if len(dlines)%2 == 1 {
		dlines = append(dlines, len(lines)-1)
	}
	//划分为以一张表为一个整体
	for i := 0; i < len(dlines); i = i + 2 {
		tl.Lines = append(tl.Lines, lines[dlines[i]:dlines[i+1]+1])
	}
	return tl
}

//tableLine2Table 表数据行变为表
func tableLine2Table(lines TableLine, classes map[string]string, dbType string, SystemName string) (tables []*Table, err error) {
	for _, tline := range lines.Lines {
		//markdown表格的表名，标题，标题数据区分行，共三行
		if len(tline) <= 3 {
			continue
		}
		var tb *Table
		for i, line := range tline {
			if i == 0 {
				//获取表名，描述名称
				name, addition, menu, err := getTableName(line)
				if err != nil {
					return nil, err
				}
				tb = NewTable(name, getTableDesc(line), getDBLink(line))
				tb.SetClass(getClass(name, classes))
				tb.Addition = addition
				tb.DBType = dbType
				tb.SystemName = SystemName
				tb.Menu = menu
				continue
			}
			if i < 3 {
				continue
			}
			c, err := line2TableColumn(line)
			if err != nil {
				return nil, err
			}
			if err := tb.AppendColumn(c); err != nil {
				return nil, err
			}
			tb.MaxSort = line.LineID
		}
		if tb != nil {
			tables = append(tables, tb)
		}
	}
	return tables, nil
}

func line2TableColumn(line *Line) (*Column, error) {
	if strings.Count(line.Text, "|") != 7 {
		return nil, fmt.Errorf("表结构有误(行:%d)", line.LineID)
	}
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return nil, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineID)
	}
	tp, l, err := getType(line)
	if err != nil {
		return nil, err
	}
	def := strings.TrimSpace(strings.Replace(colums[2], "&#124;", "|", -1))

	// if strings.Contains(tp, "char") && strings.ToLower(def) != "null" {
	// 	def = fmt.Sprintf("'%s'", def)
	// }
	c := &Column{
		Cname:  strings.TrimSpace(strings.Replace(colums[0], "&#124;", "|", -1)),
		Type:   tp,
		Len:    l,
		Def:    def,
		IsNull: strings.TrimSpace(colums[3]) == "是",
		Con:    strings.Replace(strings.TrimSpace(colums[4]), " ", "", -1),
		Desc:   strings.TrimSpace(strings.Replace(colums[5], "&#124;", "|", -1)),
		Sort:   line.LineID,
	}
	c.After = getAfter(c.Con)
	return c, nil
}

func getTableDesc(line *Line) string {
	reg := regexp.MustCompile(`[^\d^\.|\s]+[^\x00-\xff]+[^\[]+`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimSpace(names[0])
}

func getTableName(line *Line) (string, bool, string, error) {
	if !strings.HasPrefix(line.Text, "###") {
		return "", false, "", fmt.Errorf("%d行表名称标注不正确，请以###开头:%s", line.LineID, line.Text)
	}

	reg := regexp.MustCompile(`\[[\w]+[,]?[\p{Han}A-Za-z0-9_]+\]`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return "", false, "", fmt.Errorf("未设置表名称或者格式不正确:%s(行:%d)，格式：### 描述[表名,菜单名]，菜单名可选", line.Text, line.LineID)
	}
	s := strings.Split(strings.TrimRight(strings.TrimLeft(names[0], "["), "]"), ",")
	menu := ""
	if len(s) > 1 {
		menu = s[1]
	}
	return s[0], strings.Contains(line.Text, "+"), menu, nil
}

func getClass(name string, classes map[string]string) string {
	for k, v := range classes {
		if strings.HasPrefix(name, k) {
			return v
		}
	}
	return ""
}

func getDBLink(line *Line) string {
	reg := regexp.MustCompile(`\([\s\S]+\)`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimRight(strings.TrimLeft(names[0], "("), ")")
}

func getType(line *Line) (string, int, error) {
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return "", 0, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineID)
	}
	t := strings.TrimSpace(colums[1])
	reg := regexp.MustCompile(`[\w]+`)
	names := reg.FindAllString(t, -1)
	if len(names) == 0 || len(names) > 3 {
		return "", 0, fmt.Errorf("未设置字段类型:%v(行:%d)", names, line.LineID)
	}
	if len(names) == 1 {
		return t, 0, nil
	}
	l, _ := strconv.Atoi(strings.Join(names[1:], ","))
	return t, l, nil
}

func getAfter(c string) string {
	for _, v := range strings.Split(c, ",") {
		if strings.Contains(v, "After(") {
			s := strings.TrimPrefix(v, "After(")
			s = strings.TrimSuffix(s, ")")
			return s
		}
	}
	return ""
}
