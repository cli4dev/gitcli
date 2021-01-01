package tmpl

var myTypes = map[string]string{
	"date":                    "datetime",
	"decimal":                 "decimal",
	"float":                   "double",
	"int":                     "int",
	"number(1)":               "int",
	"^number\\((\\d{0,10}))$": "int",
	"^number\\((\\d{10,}))$":  "bigint",
	"timestamp":               "datetime",
	"^varchar\\(\\d+)$":       "varchar(*)",
	"^varchar2\\(\\d+)$":      "varchar(*)",
}
