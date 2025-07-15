package generator

import (
	"fmt"
	"strings"

	"github.com/orchard9/pg-goer/pkg/models"
)

type MermaidGenerator struct{}

func NewMermaidGenerator() *MermaidGenerator {
	return &MermaidGenerator{}
}

func (g *MermaidGenerator) GenerateER(schema models.Schema) (string, error) {
	var sb strings.Builder

	sb.WriteString("erDiagram\n")

	// Generate relationships first
	relationships := g.extractRelationships(schema.Tables)
	for _, rel := range relationships {
		sb.WriteString(fmt.Sprintf("    %s ||--o{ %s : %q\n", rel.ParentTable, rel.ChildTable, rel.ForeignKey))
	}

	if len(relationships) > 0 {
		sb.WriteString("\n")
	}

	// Generate table definitions
	for _, table := range schema.Tables {
		g.writeTableDefinition(&sb, &table)
	}

	return sb.String(), nil
}

type relationship struct {
	ParentTable string
	ChildTable  string
	ForeignKey  string
}

func (g *MermaidGenerator) extractRelationships(tables []models.Table) []relationship {
	var relationships []relationship

	for _, table := range tables {
		for _, fk := range table.ForeignKeys {
			// Extract referenced table name (remove schema prefix if present)
			referencedTable := fk.ReferencedTable
			if strings.Contains(referencedTable, ".") {
				parts := strings.Split(referencedTable, ".")
				referencedTable = parts[len(parts)-1]
			}

			relationships = append(relationships, relationship{
				ParentTable: referencedTable,
				ChildTable:  table.Name,
				ForeignKey:  fk.SourceColumn,
			})
		}
	}

	return relationships
}

func (g *MermaidGenerator) writeTableDefinition(sb *strings.Builder, table *models.Table) {
	fmt.Fprintf(sb, "    %s {\n", table.Name)

	for _, col := range table.Columns {
		// Normalize data type for Mermaid compatibility (single words only)
		normalizedType := g.normalizeDataType(col.DataType)
		columnDef := fmt.Sprintf("        %s %s", normalizedType, col.Name)

		// Only add the most significant constraint to follow Mermaid syntax
		// Priority: PK > UK (don't add "NOT NULL" as it's not valid Mermaid syntax)
		if col.IsPrimaryKey {
			columnDef += " PK"
		} else if col.IsUnique {
			columnDef += " UK"
		}

		sb.WriteString(columnDef + "\n")
	}

	sb.WriteString("    }\n")
}

// normalizeDataType converts PostgreSQL data types to Mermaid-compatible single words.
func (g *MermaidGenerator) normalizeDataType(dataType string) string {
	// Map of PostgreSQL types to Mermaid-friendly equivalents
	typeMap := map[string]string{
		"character varying":           "varchar",
		"timestamp without time zone": "timestamp",
		"timestamp with time zone":    "timestamptz",
		"double precision":            "float",
		"bigint":                      "bigint",
		"smallint":                    "smallint",
		"boolean":                     "boolean",
		"numeric":                     "decimal",
		"text":                        "text",
		"integer":                     "integer",
		"uuid":                        "uuid",
		"json":                        "json",
		"jsonb":                       "jsonb",
	}

	if normalized, exists := typeMap[dataType]; exists {
		return normalized
	}

	// For types not in the map, remove spaces and use first word
	parts := strings.Fields(dataType)
	if len(parts) > 0 {
		return parts[0]
	}

	return dataType
}
