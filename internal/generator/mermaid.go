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
		sb.WriteString("\n")
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
		columnDef := fmt.Sprintf("        %s %s", col.DataType, col.Name)

		var constraints []string

		if col.IsPrimaryKey {
			constraints = append(constraints, "PK")
		}

		if col.IsUnique {
			constraints = append(constraints, "UK")
		}

		if !col.IsNullable {
			constraints = append(constraints, "\"NOT NULL\"")
		}

		if len(constraints) > 0 {
			columnDef += " " + strings.Join(constraints, ",")
		}

		sb.WriteString(columnDef + "\n")
	}

	sb.WriteString("    }")
}
