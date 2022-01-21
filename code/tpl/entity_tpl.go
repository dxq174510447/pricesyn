package tpl

var EntityTpl = `package {{.EntityPackageName}}

import "time"


type {{.EntityName}} struct {
	{{range $index, $element := .Columns}}
		//{{GetAnnotation $element}}
		{{GetFieldName $element}} {{GetFieldType $element}} {{GetFieldTag $element}}
	{{end}}
}

func (m *{{.EntityName}}) TableName() string {
	return "{{.TableName}}"
}
`
