package tpl

var RepoTpl = `package {{.RepoPackageName}}

import (
	"{{.EntityPackageName}}"
)

var {{.EntityName}}Impl {{.EntityPackageName}}.{{.EntityName}} = {{.EntityPackageName}}.{{.EntityName}}{}
type {{.EntityName}}Repo struct {
	
}

`
