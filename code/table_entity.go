package code

type TableEntity struct {
	TableName         string
	EntityName        string
	EntityPackageName string
	Columns           []*TableColumnEntity
}

type TableColumnEntity struct {
	TableSchema string `gorm:"column:TABLE_SCHEMA"`

	TableEnName string `gorm:"column:TABLE_NAME"`

	ColumnName string `gorm:"primaryKey;column:COLUMN_NAME"`

	DataType string `gorm:"column:DATA_TYPE"`

	Nullable string `gorm:"column:IS_NULLABLE"`

	CharLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`

	NumericPrecision string `gorm:"column:NUMERIC_PRECISION"`

	NumericScale string `gorm:"column:NUMERIC_SCALE"`

	ColumnComment string `gorm:"column:COLUMN_COMMENT"`

	ColumnKey string `gorm:"column:COLUMN_KEY"`

	ColumnType string `gorm:"column:COLUMN_TYPE"`

	ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
}

func (m *TableColumnEntity) TableName() string {
	return "columns"
}
