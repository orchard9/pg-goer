package analyzer

import (
	"context"
	"strings"
	"testing"

	"github.com/orchard9/pg-goer/pkg/models"
)

func TestGetTables(t *testing.T) {
	tests := []struct {
		name          string
		schemas       []string
		expectTables  bool
		minTableCount int
	}{
		{
			name:          "get all public schema tables",
			schemas:       []string{"public"},
			expectTables:  true,
			minTableCount: 0,
		},
		{
			name:          "get tables from multiple schemas",
			schemas:       []string{"public", "pg_catalog"},
			expectTables:  true,
			minTableCount: 1,
		},
		{
			name:          "empty schema list returns all user tables",
			schemas:       []string{},
			expectTables:  true,
			minTableCount: 0,
		},
		{
			name:          "non-existent schema returns no tables",
			schemas:       []string{"non_existent_schema"},
			expectTables:  false,
			minTableCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skipping integration test - requires database connection")

			conn, err := Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer conn.Close()

			analyzer := NewSchemaAnalyzer(conn)
			tables, err := analyzer.GetTables(context.Background(), tt.schemas)

			if tt.expectTables && err != nil {
				t.Errorf("expected tables but got error: %v", err)
			}

			if tt.expectTables && len(tables) < tt.minTableCount {
				t.Errorf("expected at least %d tables, got %d", tt.minTableCount, len(tables))
			}

			if !tt.expectTables && len(tables) > 0 {
				t.Errorf("expected no tables but got %d", len(tables))
			}
		})
	}
}

func TestBuildTableQuery(t *testing.T) {
	analyzer := &SchemaAnalyzer{}

	tests := []struct {
		name     string
		schemas  []string
		expected string
	}{
		{
			name:     "single schema",
			schemas:  []string{"public"},
			expected: "n.nspname = $1",
		},
		{
			name:     "multiple schemas",
			schemas:  []string{"public", "custom"},
			expected: "n.nspname IN ($1, $2)",
		},
		{
			name:     "empty schemas",
			schemas:  []string{},
			expected: "NOT n.nspname IN ('pg_catalog', 'information_schema', 'pg_toast')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := analyzer.buildTableQuery(tt.schemas)
			if !contains(query, tt.expected) {
				t.Errorf("expected query to contain '%s', got: %s", tt.expected, query)
			}
		})
	}
}

func TestGetColumns(t *testing.T) {
	tests := []struct {
		name           string
		table          models.Table
		expectColumns  bool
		minColumnCount int
	}{
		{
			name: "get columns for existing table",
			table: models.Table{
				Schema: "public",
				Name:   "users",
			},
			expectColumns:  true,
			minColumnCount: 1,
		},
		{
			name: "non-existent table returns no columns",
			table: models.Table{
				Schema: "public",
				Name:   "non_existent_table",
			},
			expectColumns:  false,
			minColumnCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skipping integration test - requires database connection")

			conn, err := Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer conn.Close()

			analyzer := NewSchemaAnalyzer(conn)
			columns, err := analyzer.GetColumns(context.Background(), &tt.table)

			if tt.expectColumns && err != nil {
				t.Errorf("expected columns but got error: %v", err)
			}

			if tt.expectColumns && len(columns) < tt.minColumnCount {
				t.Errorf("expected at least %d columns, got %d", tt.minColumnCount, len(columns))
			}

			if !tt.expectColumns && len(columns) > 0 {
				t.Errorf("expected no columns but got %d", len(columns))
			}
		})
	}
}

func TestGetForeignKeys(t *testing.T) {
	tests := []struct {
		name        string
		table       models.Table
		expectKeys  bool
		minKeyCount int
	}{
		{
			name: "get foreign keys for table with relationships",
			table: models.Table{
				Schema: "public",
				Name:   "orders",
			},
			expectKeys:  true,
			minKeyCount: 0,
		},
		{
			name: "table without foreign keys returns empty",
			table: models.Table{
				Schema: "public",
				Name:   "standalone_table",
			},
			expectKeys:  false,
			minKeyCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skipping integration test - requires database connection")

			conn, err := Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer conn.Close()

			analyzer := NewSchemaAnalyzer(conn)
			foreignKeys, err := analyzer.GetForeignKeys(context.Background(), &tt.table)

			if tt.expectKeys && err != nil {
				t.Errorf("expected foreign keys but got error: %v", err)
			}

			if tt.expectKeys && len(foreignKeys) < tt.minKeyCount {
				t.Errorf("expected at least %d foreign keys, got %d", tt.minKeyCount, len(foreignKeys))
			}

			if !tt.expectKeys && len(foreignKeys) > 0 {
				t.Errorf("expected no foreign keys but got %d", len(foreignKeys))
			}
		})
	}
}

func TestGetTableRowCounts(t *testing.T) {
	tests := []struct {
		name         string
		tables       []models.Table
		expectCounts bool
	}{
		{
			name: "get row counts for existing tables",
			tables: []models.Table{
				{Schema: "public", Name: "users"},
				{Schema: "public", Name: "orders"},
			},
			expectCounts: true,
		},
		{
			name:         "empty table list returns empty map",
			tables:       []models.Table{},
			expectCounts: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skipping integration test - requires database connection")

			conn, err := Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer conn.Close()

			analyzer := NewSchemaAnalyzer(conn)
			rowCounts, err := analyzer.GetTableRowCounts(context.Background(), tt.tables)

			if tt.expectCounts && err != nil {
				t.Errorf("expected row counts but got error: %v", err)
			}

			if tt.expectCounts && len(rowCounts) == 0 && len(tt.tables) > 0 {
				t.Errorf("expected row counts but got empty map")
			}

			if !tt.expectCounts && len(rowCounts) > 0 {
				t.Errorf("expected no row counts but got %d", len(rowCounts))
			}
		})
	}
}

func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}
