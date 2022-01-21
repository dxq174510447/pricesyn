package code

import (
	"context"
	"gorm.io/gorm"
	"os"
	"pricesyn/code/tpl"
	"pricesyn/db"
	"pricesyn/util"
	"strings"
	"sync"
	"text/template"
)

var DbUser = "root"
var DbPwd = "ilove1024"
var DbHost = "10.2.246.31"
var DbPort = 3306
var DbName = "flight_order_db"
var EntityPackageName = "model"
var RepoPackageName = "repo"

type CodeFactory struct {
	db       *gorm.DB
	initLock sync.Once
}

func (c *CodeFactory) init(ctx context.Context) error {
	var err error
	c.initLock.Do(func() {
		dbfactory := &db.DbFactory{
			DbUser:     DbUser,
			DbPwd:      DbPwd,
			DbHost:     DbHost,
			DbPort:     DbPort,
			DbName:     DbName,
			DbLocation: "",
			MaxOpen:    100,
			MaxIdle:    100,
		}
		db, err1 := dbfactory.GetDb(ctx)
		if err1 != nil {
			err = err1
		} else {
			c.db = db
		}
	})
	return err
}

func (c *CodeFactory) GetTemplateMap() template.FuncMap {
	funcMap := template.FuncMap{
		"GetFieldName": func(column *TableColumnEntity) string {
			return util.StringUtil.FieldName(column.ColumnName)
		},
		"GetFieldType": func(column *TableColumnEntity) string {
			if strings.EqualFold(column.DataType, "varchar") {
				return "string"
			} else if strings.EqualFold(column.DataType, "tinyint") {
				return "int"
			} else if strings.EqualFold(column.DataType, "int") {
				return "int64"
			} else if strings.EqualFold(column.DataType, "char") {
				return "string"
			} else if strings.EqualFold(column.DataType, "datetime") {
				return "*time.Time"
			} else if strings.EqualFold(column.DataType, "date") {
				return "*time.Time"
			}
			return "string"
		},
		"GetFieldTag": func(column *TableColumnEntity) string {
			if strings.EqualFold(column.ColumnKey, "pri") {
				return "`gorm:\"primaryKey;column:" + column.ColumnName + "\"`"
			} else {
				return "`gorm:\"column:" + column.ColumnName + "\"`"
			}
		},
		"GetAnnotation": func(column *TableColumnEntity) string {
			if column.ColumnComment == "" {
				return "无"
			} else {
				return column.ColumnComment
			}
		},
	}
	return funcMap
}

func (c *CodeFactory) Generate(ctx context.Context, tablename string, schema string) error {
	table := &TableEntity{
		EntityPackageName: EntityPackageName,
		TableName:         tablename,
		EntityName:        util.StringUtil.FieldName(tablename),
	}
	columns, err := c.GetColumns(ctx, tablename, schema)
	if err != nil {
		return err
	}
	table.Columns = columns

	err = c.process(ctx, tpl.EntityTpl, table)

	return err
}

func (c *CodeFactory) process(ctx context.Context, tpl string, table *TableEntity) error {
	title := util.StringUtil.GetRandomStr(5)
	tmpl, err := template.New(title).Funcs(c.GetTemplateMap()).Parse(tpl)
	if err != nil {
		return err
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, table)
	return err
}

func (c *CodeFactory) GetColumns(ctx context.Context, tablename string, schema string) ([]*TableColumnEntity, error) {
	c.init(ctx)
	var result []*TableColumnEntity
	err := c.db.Select(`
		TABLE_SCHEMA,TABLE_NAME,COLUMN_NAME,DATA_TYPE,
		IS_NULLABLE,CHARACTER_MAXIMUM_LENGTH,NUMERIC_PRECISION,
		NUMERIC_SCALE,COLUMN_COMMENT,COLUMN_KEY,COLUMN_TYPE,COLUMN_DEFAULT
	`).Table("information_schema.columns").Where(`
		TABLE_SCHEMA = ? and TABLE_NAME= ?
	`, schema, tablename).Find(&result).Error
	if result == nil || len(result) == 0 {
		return nil, err
	}
	return result, err
}
