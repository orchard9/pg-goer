package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/orchard9/pg-goer/internal/generator"
	"github.com/orchard9/pg-goer/pkg/models"
)

type MarkdownReporter struct{}

func NewMarkdownReporter() *MarkdownReporter {
	return &MarkdownReporter{}
}

func (r *MarkdownReporter) Generate(schema models.Schema) (string, error) {
	var sb strings.Builder

	sb.WriteString("# PostgreSQL Database Documentation\n\n")
	sb.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	if len(schema.Tables) == 0 {
		sb.WriteString("No tables found in the database.\n")
		return sb.String(), nil
	}

	// Generate Table of Contents
	r.writeTableOfContents(&sb, schema.Tables)

	// Generate Database Summary
	r.writeDatabaseSummary(&sb, schema.Tables)

	// Generate Mermaid ER diagram if there are relationships
	if r.hasRelationships(schema.Tables) {
		sb.WriteString("## Database Relationships\n\n")

		mermaidGen := generator.NewMermaidGenerator()
		diagram, err := mermaidGen.GenerateER(schema)

		if err != nil {
			return "", fmt.Errorf("failed to generate Mermaid diagram: %w", err)
		}

		sb.WriteString("```mermaid\n")
		sb.WriteString(diagram)
		sb.WriteString("```\n\n")
	}

	sb.WriteString("## Tables\n\n")

	for i, table := range schema.Tables {
		if i > 0 {
			sb.WriteString("\n---\n\n")
		}

		r.writeTable(&sb, &table)
	}

	return sb.String(), nil
}

func (r *MarkdownReporter) writeTable(sb *strings.Builder, table *models.Table) {
	anchorName := strings.ToLower(strings.ReplaceAll(table.Name, "_", "-"))
	fmt.Fprintf(sb, "## %s\n\n", table.Name)
	fmt.Fprintf(sb, "<a id=\"%s\"></a>\n\n", anchorName)

	if table.Schema != "" && table.Schema != "public" {
		fmt.Fprintf(sb, "Schema: `%s`\n\n", table.Schema)
	}

	if table.RowCount > 0 {
		fmt.Fprintf(sb, "Row Count: %d\n\n", table.RowCount)
	}

	sb.WriteString("### Columns\n\n")
	sb.WriteString("| Column | Type | Nullable | Constraints | Default |\n")
	sb.WriteString("|--------|------|----------|-------------|---------|\n")

	for _, col := range table.Columns {
		r.writeColumn(sb, col)
	}
}

func (r *MarkdownReporter) writeColumn(sb *strings.Builder, col models.Column) {
	sb.WriteString("| ")
	sb.WriteString(col.Name)
	sb.WriteString(" | ")

	dataType := col.DataType
	if col.MaxLength != nil {
		dataType = fmt.Sprintf("%s(%d)", dataType, *col.MaxLength)
	}

	sb.WriteString(dataType)
	sb.WriteString(" | ")

	if col.IsNullable {
		sb.WriteString("YES")
	} else {
		sb.WriteString("NO")
	}

	sb.WriteString(" | ")

	var constraints []string
	if col.IsPrimaryKey {
		constraints = append(constraints, "PRIMARY KEY")
	}

	if col.IsUnique {
		constraints = append(constraints, "UNIQUE")
	}

	sb.WriteString(strings.Join(constraints, ", "))
	sb.WriteString(" | ")

	if col.DefaultValue != nil {
		sb.WriteString(*col.DefaultValue)
	}

	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) hasRelationships(tables []models.Table) bool {
	for _, table := range tables {
		if len(table.ForeignKeys) > 0 {
			return true
		}
	}

	return false
}

func (r *MarkdownReporter) writeTableOfContents(sb *strings.Builder, tables []models.Table) {
	sb.WriteString("## Table of Contents\n\n")

	hasRelationships := r.hasRelationships(tables)
	if hasRelationships {
		sb.WriteString("- [Database Summary](#database-summary)\n")
		sb.WriteString("- [Database Relationships](#database-relationships)\n")
	} else {
		sb.WriteString("- [Database Summary](#database-summary)\n")
	}

	sb.WriteString("- [Tables](#tables)\n")

	for _, table := range tables {
		anchorName := strings.ToLower(strings.ReplaceAll(table.Name, "_", "-"))
		sb.WriteString(fmt.Sprintf("  - [%s](#%s)\n", table.Name, anchorName))
	}

	sb.WriteString("\n")
}

func (r *MarkdownReporter) writeDatabaseSummary(sb *strings.Builder, tables []models.Table) {
	sb.WriteString("## Database Summary\n\n")

	tableCount := len(tables)
	var totalRows int64

	for _, table := range tables {
		totalRows += table.RowCount
	}

	sb.WriteString(fmt.Sprintf("**Total Tables:** %d\n", tableCount))

	if totalRows > 0 {
		sb.WriteString(fmt.Sprintf("**Total Rows:** %d\n", totalRows))
	}

	sb.WriteString("\n")
}
