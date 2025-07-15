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
	r.writeTableOfContents(&sb, schema.Tables, schema.Views, schema.Sequences, schema.Extensions)

	// Generate Database Summary
	r.writeDatabaseSummary(&sb, schema.Tables)

	// Generate Extensions section if any exist
	if len(schema.Extensions) > 0 {
		r.writeExtensions(&sb, schema.Extensions)
	}

	// Generate Views section if any exist
	if len(schema.Views) > 0 {
		r.writeViews(&sb, schema.Views)
	}

	// Generate Sequences section if any exist
	if len(schema.Sequences) > 0 {
		r.writeSequences(&sb, schema.Sequences)
	}

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

	for i := range schema.Tables {
		if i > 0 {
			sb.WriteString("\n---\n\n")
		}

		r.writeTable(&sb, &schema.Tables[i])
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

	if len(table.Indexes) > 0 {
		sb.WriteString("\n### Indexes\n\n")
		sb.WriteString("| Name | Type | Columns | Method |\n")
		sb.WriteString("|------|------|---------|--------|\n")

		for _, idx := range table.Indexes {
			r.writeIndex(sb, &idx)
		}
	}

	if len(table.Triggers) > 0 {
		sb.WriteString("\n### Triggers\n\n")
		sb.WriteString("| Name | Event | Timing | Function | Orientation |\n")
		sb.WriteString("|------|-------|--------|----------|-------------|\n")

		for _, trigger := range table.Triggers {
			r.writeTrigger(sb, &trigger)
		}
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

func (r *MarkdownReporter) writeIndex(sb *strings.Builder, idx *models.Index) {
	sb.WriteString("| ")
	sb.WriteString(idx.Name)
	sb.WriteString(" | ")
	sb.WriteString(idx.Type)
	sb.WriteString(" | ")
	sb.WriteString(strings.Join(idx.Columns, ", "))
	sb.WriteString(" | ")
	sb.WriteString(idx.Method)
	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) writeTrigger(sb *strings.Builder, trigger *models.Trigger) {
	sb.WriteString("| ")
	sb.WriteString(trigger.Name)
	sb.WriteString(" | ")
	sb.WriteString(trigger.Event)
	sb.WriteString(" | ")
	sb.WriteString(trigger.Timing)
	sb.WriteString(" | ")
	sb.WriteString(trigger.Function)
	sb.WriteString(" | ")
	sb.WriteString(trigger.Orientation)
	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) hasRelationships(tables []models.Table) bool {
	for i := range tables {
		if len(tables[i].ForeignKeys) > 0 {
			return true
		}
	}

	return false
}

func (r *MarkdownReporter) writeTableOfContents(sb *strings.Builder, tables []models.Table, views []models.View, sequences []models.Sequence, extensions []models.Extension) {
	sb.WriteString("## Table of Contents\n\n")

	sb.WriteString("- [Database Summary](#database-summary)\n")

	if len(extensions) > 0 {
		sb.WriteString("- [PostgreSQL Extensions](#postgresql-extensions)\n")
	}

	if len(views) > 0 {
		sb.WriteString("- [Views](#views)\n")
	}

	if len(sequences) > 0 {
		sb.WriteString("- [Sequences](#sequences)\n")
	}

	hasRelationships := r.hasRelationships(tables)
	if hasRelationships {
		sb.WriteString("- [Database Relationships](#database-relationships)\n")
	}

	sb.WriteString("- [Tables](#tables)\n")

	for i := range tables {
		anchorName := strings.ToLower(strings.ReplaceAll(tables[i].Name, "_", "-"))
		fmt.Fprintf(sb, "  - [%s](#%s)\n", tables[i].Name, anchorName)
	}

	sb.WriteString("\n")
}

func (r *MarkdownReporter) writeExtensions(sb *strings.Builder, extensions []models.Extension) {
	sb.WriteString("## PostgreSQL Extensions\n\n")

	if len(extensions) == 0 {
		sb.WriteString("No extensions are installed.\n\n")
		return
	}

	sb.WriteString("| Extension | Version | Schema |\n")
	sb.WriteString("|-----------|---------|--------|\n")

	for i := range extensions {
		r.writeExtension(sb, &extensions[i])
	}

	sb.WriteString("\n")
}

func (r *MarkdownReporter) writeExtension(sb *strings.Builder, ext *models.Extension) {
	sb.WriteString("| ")
	sb.WriteString(ext.Name)
	sb.WriteString(" | ")
	sb.WriteString(ext.Version)
	sb.WriteString(" | ")
	sb.WriteString(ext.Schema)
	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) writeViews(sb *strings.Builder, views []models.View) {
	sb.WriteString("## Views\n\n")

	if len(views) == 0 {
		sb.WriteString("No views are defined.\n\n")
		return
	}

	sb.WriteString("| View | Schema |\n")
	sb.WriteString("|------|--------|\n")

	for i := range views {
		r.writeView(sb, &views[i])
	}

	sb.WriteString("\n")
}

func (r *MarkdownReporter) writeView(sb *strings.Builder, view *models.View) {
	sb.WriteString("| ")
	sb.WriteString(view.Name)
	sb.WriteString(" | ")
	sb.WriteString(view.Schema)
	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) writeSequences(sb *strings.Builder, sequences []models.Sequence) {
	sb.WriteString("## Sequences\n\n")

	if len(sequences) == 0 {
		sb.WriteString("No sequences are defined.\n\n")
		return
	}

	sb.WriteString("| Sequence | Schema | Data Type | Start | Min | Max | Increment |\n")
	sb.WriteString("|----------|--------|-----------|-------|-----|-----|----------|\n")

	for i := range sequences {
		r.writeSequence(sb, &sequences[i])
	}

	sb.WriteString("\n")
}

func (r *MarkdownReporter) writeSequence(sb *strings.Builder, seq *models.Sequence) {
	sb.WriteString("| ")
	sb.WriteString(seq.Name)
	sb.WriteString(" | ")
	sb.WriteString(seq.Schema)
	sb.WriteString(" | ")
	sb.WriteString(seq.DataType)
	sb.WriteString(" | ")
	fmt.Fprintf(sb, "%d", seq.StartValue)
	sb.WriteString(" | ")
	fmt.Fprintf(sb, "%d", seq.MinValue)
	sb.WriteString(" | ")
	fmt.Fprintf(sb, "%d", seq.MaxValue)
	sb.WriteString(" | ")
	fmt.Fprintf(sb, "%d", seq.Increment)
	sb.WriteString(" |\n")
}

func (r *MarkdownReporter) writeDatabaseSummary(sb *strings.Builder, tables []models.Table) {
	sb.WriteString("## Database Summary\n\n")

	tableCount := len(tables)

	var totalRows int64

	for i := range tables {
		totalRows += tables[i].RowCount
	}

	fmt.Fprintf(sb, "**Total Tables:** %d\n", tableCount)

	if totalRows > 0 {
		fmt.Fprintf(sb, "**Total Rows:** %d\n", totalRows)
	}

	sb.WriteString("\n")
}
