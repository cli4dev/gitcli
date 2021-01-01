package tmpts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

var dataType map[string]string
var dataTypeScope map[string]map[string]string

//参考资料 http://www.sqlines.com/oracle-to-mysql
func init() {
	dataType = make(map[string]string)
	dataType = map[string]string{
		"BFILE":                    "VARCHAR(255)",
		"BINARY_FLOAT":             "FLOAT",
		"BINARY_DOUBLE":            "DOUBLE",
		"BLOB":                     "LONGBLOB",
		"CHAR":                     "*",
		"CHARACTER":                "*",
		"CLOB":                     "LONGTEXT",
		"DATE":                     "DATETIME",
		"DECIMAL":                  "DECIMAL",
		"DOUBLE PRECISION":         "DOUBLE PRECISION",
		"FLOAT":                    "DOUBLE",
		"INTEGER":                  "INT", //或者	DECIMAL(38)
		"INT":                      "INT", //或者   DECIMAL(38)
		"INTERVAL YEAR TO MONTH":   "VARCHAR(30)",
		"INTERVAL DAY TO SECOND":   "VARCHAR(30)",
		"LONG":                     "LONGTEXT",
		"LONG RAW":                 "LONGBLOB",
		"NCHAR":                    "*",
		"NCHAR VARYING":            "NCHAR VARYING",
		"NCLOB":                    "NVARCHAR(4000)",
		"NUMBER":                   "*",
		"NUMERIC":                  "NUMERIC",
		"NVARCHAR2":                "NVARCHAR",
		"RAW":                      "*",
		"REAL":                     "DOUBLE",
		"ROWID":                    "CHAR(10)",
		"SMALLINT":                 "DECIMAL(38)",
		"TIMESTAMP":                "DATETIME",
		"TIMESTAMP WITH TIME ZONE": "DATETIME",
		"UROWID":                   "VARCHAR",
		"VARCHAR":                  "VARCHAR",
		"VARCHAR2":                 "VARCHAR",
		"XMLTYPE":                  "LONGTEXT",
	}

	//范围为左闭右开
	dataTypeScope = make(map[string]map[string]string)
	dataTypeScope = map[string]map[string]string{
		"CHAR": map[string]string{
			"1~256":    "CHAR",
			"256~2000": "VARCHAR",
		},
		"CHARACTER": map[string]string{
			"1~256":    "CHARACTER",
			"256~2000": "VARCHAR",
		},
		"NCHAR": map[string]string{
			"1~256":    "NCHAR",
			"256~2000": "NVARCHAR",
		},
		"NUMBER": map[string]string{ //NUMBER(p,0), NUMBER(p),NUMBER(p,s),NUMBER, NUMBER(*)
			"1~3":   "TINYINT",  //p>0
			"3~5":   "SMALLINT", //p>0
			"5~9":   "INT",      //p>0
			"9~19":  "BIGINT",   //p>0
			"19~38": "BIGINT",   //p>0,s=0 理论上是转为DECIMAL，因为超过了mysql最大的bigint；但是业务中主键使用了number(20)，转化为biginit才能保证语句执行正确
			"s~s":   "DECIMAL",  //s>0
			"-~-":   "DOUBLE",   //NUMBER
			"*~*":   "DOUBLE",   //NUMBER(*)
		},
		"RAW": map[string]string{
			"1~256":    "BINARY",
			"256~2000": "VARBINARY",
		},
	}
}

//ConvertDataType .
func ConvertDataType(oracleType string) (mysqlType string, err error) {
	reg := regexp.MustCompile(`[\w]+`)
	tps := reg.FindAllString(strings.ToUpper(oracleType), -1)
	typeLen := len(tps)
	if typeLen == 0 || typeLen > 3 {
		return "", fmt.Errorf("数据类型%s不能进行解析", oracleType)
	}
	t := tps[0]
	if _, ok := dataType[t]; !ok {
		return "", fmt.Errorf("未匹配到%s对应的mysql数据类型", oracleType)
	}

	mysqlType = dataType[t]
	if mysqlType == "*" {
		if tps[0] == "NUMBER" {
			mysqlType, err = getNumDataTypeScope(tps, typeLen)
			if err != nil {
				return "", fmt.Errorf("类型%s转换错误:%+v", oracleType, err)
			}
		} else {
			mysqlType, err = getdataTypeScope(tps[0], types.GetInt(tps[1], -1))
			if err != nil {
				return "", fmt.Errorf("类型%s转换错误:%+v", oracleType, err)
			}
		}
	}

	if typeLen == 2 {
		if index := strings.Index(oracleType, "("); index >= 0 {
			return fmt.Sprintf("%s(%s)", oracleType[:index], tps[1]), nil
		}
		return fmt.Sprintf("%s(%s)", mysqlType, tps[1]), nil
	}
	if typeLen == 3 {
		if index := strings.Index(oracleType, "("); index >= 0 {
			return fmt.Sprintf("%s(%s,%s)", oracleType[:index], tps[1], tps[2]), nil
		}
		return fmt.Sprintf("%s(%s,%s)", mysqlType, tps[1], tps[2]), nil
	}
	return mysqlType, nil
}

func getNumDataTypeScope(tps []string, len int) (string, error) {
	if len == 1 {
		return dataTypeScope[tps[0]]["-~-"], nil
	}
	if len == 2 {
		if tps[1] == "*" {
			return dataTypeScope[tps[0]]["*~*"], nil
		}
		return getdataTypeScope(tps[0], types.GetInt(tps[1], -1))
	}
	if len == 3 {
		if types.GetInt(tps[2], -1) > 0 {
			return dataTypeScope[tps[0]]["s~s"], nil
		}
		if types.GetInt(tps[2], -1) == 0 {
			return getdataTypeScope(tps[0], types.GetInt(tps[1], -1))
		}
	}
	return "", fmt.Errorf("未找到NUMBER类型mysql的范围对应长度")
}

func getdataTypeScope(t string, size int) (string, error) {
	if _, ok := dataTypeScope[t]; !ok {
		return "", fmt.Errorf("未找到%s的mysql的相应%d长度的数据类型", t, size)
	}
	for k, v := range dataTypeScope[t] {
		scope := strings.Split(k, "~")
		if types.GetInt(scope[0], -1) <= size && size < types.GetInt(scope[1], -1) {
			return v, nil
		}
	}
	return "", fmt.Errorf("未设置%s的mysql的相应%d长度的数据类型", t, size)
}
