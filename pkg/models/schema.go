package models

type Schema struct {
	Name   string
	Tables []Table
}

type Table struct {
	Schema      string
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
	RowCount    int64
}

type Column struct {
	Name         string
	DataType     string
	IsNullable   bool
	DefaultValue *string
	IsPrimaryKey bool
	IsUnique     bool
	MaxLength    *int
}

type ForeignKey struct {
	Name             string
	SourceTable      string
	SourceColumn     string
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string
	OnUpdate         string
}
