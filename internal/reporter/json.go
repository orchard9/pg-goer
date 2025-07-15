package reporter

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/orchard9/pg-goer/pkg/models"
)

type JSONReporter struct{}

func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

// JSONOutput represents the JSON structure for database documentation.
type JSONOutput struct {
	GeneratedAt   string             `json:"generated_at"`
	DatabaseName  string             `json:"database_name"`
	Summary       DatabaseSummary    `json:"summary"`
	Extensions    []JSONExtension    `json:"extensions,omitempty"`
	Tables        []JSONTable        `json:"tables"`
	Relationships []JSONRelationship `json:"relationships,omitempty"`
}

type DatabaseSummary struct {
	TableCount int   `json:"table_count"`
	TotalRows  int64 `json:"total_rows"`
}

type JSONTable struct {
	Name        string           `json:"name"`
	Schema      string           `json:"schema"`
	RowCount    int64            `json:"row_count"`
	Columns     []JSONColumn     `json:"columns"`
	ForeignKeys []JSONForeignKey `json:"foreign_keys,omitempty"`
	Indexes     []JSONIndex      `json:"indexes,omitempty"`
	Triggers    []JSONTrigger    `json:"triggers,omitempty"`
}

type JSONColumn struct {
	Name         string  `json:"name"`
	DataType     string  `json:"data_type"`
	MaxLength    *int    `json:"max_length,omitempty"`
	IsNullable   bool    `json:"is_nullable"`
	IsPrimaryKey bool    `json:"is_primary_key"`
	IsUnique     bool    `json:"is_unique"`
	DefaultValue *string `json:"default_value,omitempty"`
}

type JSONForeignKey struct {
	Name             string `json:"name"`
	SourceTable      string `json:"source_table"`
	SourceColumn     string `json:"source_column"`
	ReferencedTable  string `json:"referenced_table"`
	ReferencedColumn string `json:"referenced_column"`
}

type JSONIndex struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	IsPrimary bool     `json:"is_primary"`
	IsUnique  bool     `json:"is_unique"`
	Columns   []string `json:"columns"`
	Method    string   `json:"method"`
}

type JSONTrigger struct {
	Name        string `json:"name"`
	Event       string `json:"event"`
	Timing      string `json:"timing"`
	Function    string `json:"function"`
	Orientation string `json:"orientation"`
}

type JSONExtension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  string `json:"schema"`
}

type JSONRelationship struct {
	ParentTable string `json:"parent_table"`
	ChildTable  string `json:"child_table"`
	ForeignKey  string `json:"foreign_key"`
}

func (r *JSONReporter) Generate(schema *models.Schema) (string, error) {
	output := JSONOutput{
		GeneratedAt:   time.Now().Format(time.RFC3339),
		DatabaseName:  schema.Name,
		Summary:       r.buildSummary(schema.Tables),
		Extensions:    r.buildExtensions(schema.Extensions),
		Tables:        r.buildTables(schema.Tables),
		Relationships: r.buildRelationships(schema.Tables),
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (r *JSONReporter) buildSummary(tables []models.Table) DatabaseSummary {
	var totalRows int64
	for i := range tables {
		totalRows += tables[i].RowCount
	}

	return DatabaseSummary{
		TableCount: len(tables),
		TotalRows:  totalRows,
	}
}

func (r *JSONReporter) buildTables(tables []models.Table) []JSONTable {
	jsonTables := make([]JSONTable, len(tables))

	for i := range tables {
		table := &tables[i]
		jsonTables[i] = JSONTable{
			Name:        table.Name,
			Schema:      table.Schema,
			RowCount:    table.RowCount,
			Columns:     r.buildColumns(table.Columns),
			ForeignKeys: r.buildForeignKeys(table.ForeignKeys),
			Indexes:     r.buildIndexes(table.Indexes),
			Triggers:    r.buildTriggers(table.Triggers),
		}
	}

	return jsonTables
}

func (r *JSONReporter) buildColumns(columns []models.Column) []JSONColumn {
	jsonColumns := make([]JSONColumn, len(columns))

	for i, col := range columns {
		jsonColumns[i] = JSONColumn{
			Name:         col.Name,
			DataType:     col.DataType,
			MaxLength:    col.MaxLength,
			IsNullable:   col.IsNullable,
			IsPrimaryKey: col.IsPrimaryKey,
			IsUnique:     col.IsUnique,
			DefaultValue: col.DefaultValue,
		}
	}

	return jsonColumns
}

func (r *JSONReporter) buildForeignKeys(foreignKeys []models.ForeignKey) []JSONForeignKey {
	if len(foreignKeys) == 0 {
		return nil
	}

	jsonForeignKeys := make([]JSONForeignKey, len(foreignKeys))

	for i, fk := range foreignKeys {
		jsonForeignKeys[i] = JSONForeignKey{
			Name:             fk.Name,
			SourceTable:      fk.SourceTable,
			SourceColumn:     fk.SourceColumn,
			ReferencedTable:  fk.ReferencedTable,
			ReferencedColumn: fk.ReferencedColumn,
		}
	}

	return jsonForeignKeys
}

func (r *JSONReporter) buildIndexes(indexes []models.Index) []JSONIndex {
	if len(indexes) == 0 {
		return nil
	}

	jsonIndexes := make([]JSONIndex, len(indexes))

	for i, idx := range indexes {
		jsonIndexes[i] = JSONIndex{
			Name:      idx.Name,
			Type:      idx.Type,
			IsPrimary: idx.IsPrimary,
			IsUnique:  idx.IsUnique,
			Columns:   idx.Columns,
			Method:    idx.Method,
		}
	}

	return jsonIndexes
}

func (r *JSONReporter) buildTriggers(triggers []models.Trigger) []JSONTrigger {
	if len(triggers) == 0 {
		return nil
	}

	jsonTriggers := make([]JSONTrigger, len(triggers))

	for i, trigger := range triggers {
		jsonTriggers[i] = JSONTrigger{
			Name:        trigger.Name,
			Event:       trigger.Event,
			Timing:      trigger.Timing,
			Function:    trigger.Function,
			Orientation: trigger.Orientation,
		}
	}

	return jsonTriggers
}

func (r *JSONReporter) buildExtensions(extensions []models.Extension) []JSONExtension {
	if len(extensions) == 0 {
		return nil
	}

	jsonExtensions := make([]JSONExtension, len(extensions))

	for i, ext := range extensions {
		jsonExtensions[i] = JSONExtension{
			Name:    ext.Name,
			Version: ext.Version,
			Schema:  ext.Schema,
		}
	}

	return jsonExtensions
}

func (r *JSONReporter) buildRelationships(tables []models.Table) []JSONRelationship {
	var relationships []JSONRelationship

	for i := range tables {
		for _, fk := range tables[i].ForeignKeys {
			// Extract referenced table name (remove schema prefix if present)
			referencedTable := fk.ReferencedTable
			if referencedTable != "" && referencedTable[0] != '"' {
				parts := strings.Split(referencedTable, ".")
				if len(parts) > 1 {
					referencedTable = parts[len(parts)-1]
				}
			}

			relationships = append(relationships, JSONRelationship{
				ParentTable: referencedTable,
				ChildTable:  tables[i].Name,
				ForeignKey:  fk.SourceColumn,
			})
		}
	}

	return relationships
}
