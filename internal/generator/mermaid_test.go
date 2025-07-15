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
		})
	}
}
