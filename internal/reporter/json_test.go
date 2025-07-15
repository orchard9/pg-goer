package reporter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/orchard9/pg-goer/pkg/models"
)

func TestJSONReporter_Generate(t *testing.T) {
	tests := []struct {
		name           string
		schema         models.Schema
		expectContains []string
		expectFields   []string
	}{
		{
			name: "basic schema with tables",
			schema: models.Schema{
				Name: "test_db",
				Tables: []models.Table{
					{
						Schema:   "public",
						Name:     "users",
						RowCount: 10,
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true, IsNullable: false},
							{Name: "email", DataType: "varchar", IsUnique: true, IsNullable: false},
							{Name: "name", DataType: "varchar", IsNullable: true},
						},
					},
					{
						Schema:   "public",
						Name:     "posts",
						RowCount: 25,
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true, IsNullable: false},
							{Name: "user_id", DataType: "integer", IsNullable: false},
							{Name: "title", DataType: "varchar", IsNullable: false},
						},
						ForeignKeys: []models.ForeignKey{
							{
								Name:             "fk_posts_user_id",
								SourceTable:      "posts",
								SourceColumn:     "user_id",
								ReferencedTable:  "public.users",
								ReferencedColumn: "id",
							},
						},
					},
				},
			},
			expectContains: []string{
				`"database_name": "test_db"`,
				`"table_count": 2`,
				`"total_rows": 35`,
				`"name": "users"`,
				`"name": "posts"`,
				`"is_primary_key": true`,
				`"is_unique": true`,
				`"foreign_keys"`,
				`"relationships"`,
			},
			expectFields: []string{
				"generated_at",
				"database_name",
				"summary",
				"tables",
				"relationships",
			},
		},
		{
			name: "schema with no relationships",
			schema: models.Schema{
				Name: "simple_db",
				Tables: []models.Table{
					{
						Schema:   "public",
						Name:     "products",
						RowCount: 5,
						Columns: []models.Column{
							{Name: "id", DataType: "integer", IsPrimaryKey: true, IsNullable: false},
							{Name: "name", DataType: "varchar", IsNullable: false},
						},
					},
				},
			},
			expectContains: []string{
				`"database_name": "simple_db"`,
				`"table_count": 1`,
				`"total_rows": 5`,
				`"name": "products"`,
			},
			expectFields: []string{
				"generated_at",
				"database_name",
				"summary",
				"tables",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewJSONReporter()
			output, err := reporter.Generate(tt.schema)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify it's valid JSON
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			// Check that all expected fields are present
			for _, field := range tt.expectFields {
				if _, exists := jsonData[field]; !exists {
					t.Errorf("expected field '%s' not found in JSON output", field)
				}
			}

			// Check that output contains expected strings
			for _, expected := range tt.expectContains {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain '%s', but it didn't.\nOutput:\n%s", expected, output)
				}
			}

			// Verify JSON structure more deeply
			validateJSONStructure(t, jsonData, tt.schema)
		})
	}
}

func validateJSONStructure(t *testing.T, jsonData map[string]interface{}, schema models.Schema) {
	t.Helper()

	// Check summary
	summary, ok := jsonData["summary"].(map[string]interface{})
	if !ok {
		t.Error("summary should be an object")
		return
	}

	expectedTableCount := float64(len(schema.Tables))
	if tableCount, ok := summary["table_count"].(float64); !ok || tableCount != expectedTableCount {
		t.Errorf("expected table_count to be %v, got %v", expectedTableCount, summary["table_count"])
	}

	// Check tables array
	tables, ok := jsonData["tables"].([]interface{})
	if !ok {
		t.Error("tables should be an array")
		return
	}

	if len(tables) != len(schema.Tables) {
		t.Errorf("expected %d tables, got %d", len(schema.Tables), len(tables))
	}

	// Check first table structure if it exists
	if len(tables) > 0 {
		firstTable, ok := tables[0].(map[string]interface{})
		if !ok {
			t.Error("table should be an object")
			return
		}

		requiredTableFields := []string{"name", "schema", "row_count", "columns"}
		for _, field := range requiredTableFields {
			if _, exists := firstTable[field]; !exists {
				t.Errorf("table should have field '%s'", field)
			}
		}

		// Check columns array
		columns, ok := firstTable["columns"].([]interface{})
		if !ok {
			t.Error("columns should be an array")
			return
		}

		if len(columns) > 0 {
			firstColumn, ok := columns[0].(map[string]interface{})
			if !ok {
				t.Error("column should be an object")
				return
			}

			requiredColumnFields := []string{"name", "data_type", "is_nullable", "is_primary_key", "is_unique"}
			for _, field := range requiredColumnFields {
				if _, exists := firstColumn[field]; !exists {
					t.Errorf("column should have field '%s'", field)
				}
			}
		}
	}
}

func TestJSONReporter_EmptySchema(t *testing.T) {
	reporter := NewJSONReporter()
	schema := models.Schema{
		Name:   "empty_db",
		Tables: []models.Table{},
	}

	output, err := reporter.Generate(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Check that tables array is empty
	tables, ok := jsonData["tables"].([]interface{})
	if !ok {
		t.Error("tables should be an array")
		return
	}

	if len(tables) != 0 {
		t.Errorf("expected empty tables array, got %d tables", len(tables))
	}

	// Check that summary shows 0 tables
	summary, ok := jsonData["summary"].(map[string]interface{})
	if !ok {
		t.Error("summary should be an object")
		return
	}

	if tableCount, ok := summary["table_count"].(float64); !ok || tableCount != 0 {
		t.Errorf("expected table_count to be 0, got %v", summary["table_count"])
	}
}
