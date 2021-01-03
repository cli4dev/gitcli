package tmpl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

//Tables markdown中的所有表信息
type Tables struct {
	PKG string
	Tbs []*Table
}

//SetPkg 添加行信息
func (t *Tables) SetPkg(path string) {
	names := strings.Split(strings.Trim(path, "/"), "/")
	t.PKG = names[len(names)-1]
	for _, tb := range t.Tbs {
		tb.SetPkg(t.PKG)
	}
}

//Line 每一行信息
type Line struct {
	Text   string
	LineID int
}

//TableLine 表的每一行
type TableLine struct {
	Lines [][]*Line
}

//Markdown2DB 读取markdown文件并转换为MarkDownDB对象
func Markdown2DB(fn string) (*Tables, error) {
	lines, err := readMarkdown(fn)
	if err != nil {
		return nil, err
	}

	return tableLine2Table(line2TableLine(lines))
}

//readMarkdown 读取md文件
func readMarkdown(name string) ([]*Line, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readMarkdownByReader(bufio.NewReader(f))
}

//readMarkdown 读取md文件
func readMarkdownByReader(rd *bufio.Reader) ([]*Line, error) {
	lines := make([]*Line, 0, 64)
	num := 0
	for {
		num++
		line, err := rd.ReadString('\n')
		if line == "" && (err != nil || io.EOF == err) {
			break
		}
		line = strings.Trim(strings.Trim(line, "\n"), "\t")
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, &Line{Text: line, LineID: num})
	}
	return lines, nil
}

//lines2TableLine 数据行转变为以表为单个整体的数据行
func line2TableLine(lines []*Line) (tl TableLine) {
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
func tableLine2Table(lines TableLine) (tables *Tables, err error) {
	tables = &Tables{Tbs: make([]*Table, 0, 1)}
	for _, tline := range lines.Lines {
		//markdown表格的表名，标题，标题数据区分行，共三行
		if len(tline) <= 3 {
			continue
		}
		var tb *Table
		for i, line := range tline {
			if i == 0 {
				//获取表名，描述名称
				name, err := getTableName(line)
				if err != nil {
					return nil, err
				}
				tb = NewTable(name, getTableDesc(line))
				continue
			}
			if i < 3 {
				continue
			}
			c, err := line2TableRow(line)
			if err != nil {
				return nil, err
			}
			if err := tb.AddRow(c); err != nil {
				return nil, err
			}
		}
		if tb != nil {
			tables.Tbs = append(tables.Tbs, tb)
		}
	}
	return tables, nil
}

func line2TableRow(line *Line) (*Row, error) {
	if strings.Count(line.Text, "|") != 7 {
		return nil, fmt.Errorf("表结构有误(行:%d)", line.LineID)
	}
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return nil, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineID)
	}
	c := &Row{
		Name:   strings.TrimSpace(strings.Replace(colums[0], "&#124;", "|", -1)),
		Type:   strings.TrimSpace(colums[1]),
		Def:    strings.TrimSpace(strings.Replace(colums[2], "&#124;", "|", -1)),
		IsNull: strings.TrimSpace(colums[3]),
		Con:    strings.Replace(strings.TrimSpace(colums[4]), " ", "", -1),
		Desc:   strings.TrimSpace(strings.Replace(colums[5], "&#124;", "|", -1)),
	}
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

func getTableName(line *Line) (string, error) {
	if !strings.HasPrefix(line.Text, "###") {
		return "", fmt.Errorf("%d行表名称标注不正确，请以###开头:(%s)", line.LineID, line.Text)
	}

	reg := regexp.MustCompile(`\[[\w]+[,]?[\p{Han}A-Za-z0-9_]+\]`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return "", fmt.Errorf("未设置表名称或者格式不正确:%s(行:%d)，格式：### 描述[表名,菜单名]，菜单名可选", line.Text, line.LineID)
	}
	s := strings.Split(strings.TrimRight(strings.TrimLeft(names[0], "["), "]"), ",")
	return s[0], nil
}