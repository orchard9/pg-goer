package reporter

import (
	"strings"
	"testing"

	"github.com/orchard9/pg-goer/pkg/models"
)

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name           string
		schema         models.Schema
		expectContains []string
	}{
		{
			name: "basic table documentation",
			schema: models.Schema{
				Name: "public",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "users",
						Columns: []models.Column{
							{
								Name:         "id",
								DataType:     "integer",
								IsNullable:   false,
								IsPrimaryKey: true,
							},
							{
								Name:       "email",
								DataType:   "varchar",
								IsNullable: false,
								MaxLength:  intPtr(255),
							},
						},
					},
				},
			},
			expectContains: []string{
				"# PostgreSQL Database Documentation",
				"## users",
				"| Column | Type | Nullable | Constraints |",
				"| id | integer | NO | PRIMARY KEY |",
				"| email | varchar(255) | NO |",
			},
		},
		{
			name: "table with default values",
			schema: models.Schema{
				Name: "public",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "config",
						Columns: []models.Column{
							{
								Name:         "key",
								DataType:     "varchar",
								IsNullable:   false,
								DefaultValue: stringPtr("''::character varying"),
							},
							{
								Name:         "value",
								DataType:     "text",
								IsNullable:   true,
								DefaultValue: stringPtr("NULL"),
							},
						},
					},
				},
			},
			expectContains: []string{
				"## config",
				"| key | varchar | NO |  | ''::character varying |",
				"| value | text | YES |  | NULL |",
			},
		},
		{
			name: "schema with Mermaid diagram",
			schema: models.Schema{
				Name: "public",
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
				"# PostgreSQL Database Documentation",
				"## Database Relationships",
				"```mermaid",
				"erDiagram",
				"users ||--o{ orders",
				"```",
				"## users",
				"## orders",
			},
		},
		{
			name: "complete formatted report with TOC",
			schema: models.Schema{
				Name: "sample_db",
				Tables: []models.Table{
					{
						Schema: "public",
						Name:   "users",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "email", DataType: "varchar", MaxLength: intPtr(255)},
						},
						RowCount: 1500,
					},
					{
						Schema: "public",
						Name:   "orders",
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true},
							{Name: "user_id", DataType: "integer"},
							{Name: "total", DataType: "decimal"},
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
						RowCount: 8500,
					},
				},
			},
			expectContains: []string{
				"# PostgreSQL Database Documentation",
				"## Table of Contents",
				"- [Database Relationships](#database-relationships)",
				"- [Tables](#tables)",
				"  - [users](#users)",
				"  - [orders](#orders)",
				"## Database Summary",
				"**Total Tables:** 2",
				"**Total Rows:** 10000",
				"Row Count: 1500",
				"Row Count: 8500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewMarkdownReporter()
			output, err := reporter.Generate(tt.schema)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, expected := range tt.expectContains {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain '%s', but it didn't.\nOutput:\n%s", expected, output)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
