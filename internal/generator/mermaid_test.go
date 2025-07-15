package generator

import (
	"strings"
	"testing"

	"github.com/orchard9/pg-goer/pkg/models"
)

func TestGenerateMermaidER(t *testing.T) {
	tests := []struct {
		name           string
		schema         models.Schema
		expectContains []string
	}{
		{
			name: "basic tables without relationships",
			schema: models.Schema{
				Name: "test_db",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "users",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "email", DataType: "varchar"},
						},
					},
					{
						Schema: "public",
						Name:   "posts",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "title", DataType: "varchar"},
						},
					},
				},
			},
			expectContains: []string{
				"erDiagram",
				"users {",
				"integer id PK",
				"varchar email",
				"posts {",
				"integer id PK",
				"varchar title",
			},
		},
		{
			name: "tables with foreign key relationships",
			schema: models.Schema{
				Name: "test_db",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "users",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
						},
					},
					{
						Schema: "public",
						Name:   "orders",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "user_id", DataType: "integer"},
						},
						ForeignKeys: []models.ForeignKey{
							{
								Name:             "fk_orders_user_id",
								SourceTable:      "orders",
								SourceColumn:     "user_id",
								ReferencedTable:  "public.users",
								ReferencedColumn: "id",
							},
						},
					},
				},
			},
			expectContains: []string{
				"erDiagram",
				"users ||--o{ orders : \"user_id\"",
				"users {",
				"orders {",
			},
		},
		{
			name: "comprehensive constraints and syntax validation",
			schema: models.Schema{
				Name: "test_db",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "products",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true, IsNullable: false},
							{Name: "code", DataType: "varchar", IsUnique: true, IsNullable: false},
							{Name: "name", DataType: "varchar", IsNullable: false},
							{Name: "description", DataType: "text", IsNullable: true},
							{Name: "price", DataType: "decimal", IsNullable: true},
						},
					},
				},
			},
			expectContains: []string{
				"erDiagram",
				"products {",
				"        integer id PK",
				"        varchar code UK",
				"        varchar name",
				"        text description",
				"        decimal price",
			},
		},
		{
			name: "data types with spaces should be normalized",
			schema: models.Schema{
				Name: "test_db",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "test_table",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "name", DataType: "character varying", IsUnique: true},
							{Name: "created_at", DataType: "timestamp without time zone"},
							{Name: "updated_at", DataType: "timestamp with time zone"},
							{Name: "amount", DataType: "double precision"},
						},
					},
				},
			},
			expectContains: []string{
				"erDiagram",
				"test_table {",
				"        integer id PK",
				"        varchar name UK",
				"        timestamp created_at",
				"        timestamptz updated_at",
				"        float amount",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewMermaidGenerator()
			output, err := generator.GenerateER(tt.schema)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, expected := range tt.expectContains {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain '%s', but it didn't.\nOutput:\n%s", expected, output)
				}
			}

			// Validate that output doesn't contain invalid Mermaid syntax
			validateMermaidSyntax(t, output)
		})
	}
}

// validateMermaidSyntax checks for common Mermaid syntax errors that would cause parser failures.
func validateMermaidSyntax(t *testing.T, mermaidOutput string) {
	t.Helper()

	// Check for invalid patterns that cause parse errors
	invalidPatterns := []string{
		",\"",          // Comma followed by quote (like PK,"NOT NULL")
		"PK,",          // Primary key followed by comma
		"UK,",          // Unique key followed by comma
		"\"NOT NULL\"", // Quoted NOT NULL (not valid in Mermaid)
	}

	for _, pattern := range invalidPatterns {
		if strings.Contains(mermaidOutput, pattern) {
			t.Errorf("Mermaid output contains invalid syntax pattern '%s'.\nOutput:\n%s", pattern, mermaidOutput)
		}
	}

	// Validate basic structure
	if !strings.Contains(mermaidOutput, "erDiagram") {
		t.Error("Mermaid output should start with 'erDiagram'")
	}

	// Check that all table definitions have proper closing braces
	// Count only table definition braces (lines that end with " {")
	lines := strings.Split(mermaidOutput, "\n")
	tableOpenBraces := 0
	closeBraces := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, " {") {
			tableOpenBraces++
		}

		if trimmed == "}" {
			closeBraces++
		}
	}

	if tableOpenBraces != closeBraces {
		t.Errorf("Unmatched table braces in Mermaid output: %d table opens, %d closes.\nOutput:\n%s", tableOpenBraces, closeBraces, mermaidOutput)
	}
}
