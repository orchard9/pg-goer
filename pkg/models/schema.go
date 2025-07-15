package models

type Schema struct {
	Name       string
	Tables     []Table
	Views      []View
	Sequences  []Sequence
	Extensions []Extension
}

type Table struct {
	Schema      string
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
	Indexes     []Index
	Triggers    []Trigger
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

type Index struct {
	Name      string
	Type      string
	IsPrimary bool
	IsUnique  bool
	Columns   []string
	Method    string
}

type Trigger struct {
	Name        string
	Event       string
	Timing      string
	Function    string
	Orientation string
}

type Extension struct {
	Name    string
	Version string
	Schema  string
}

type View struct {
	Schema string
	Name   string
}

type Sequence struct {
	Schema     string
	Name       string
	DataType   string
	StartValue int64
	MinValue   int64
	MaxValue   int64
	Increment  int64
}
