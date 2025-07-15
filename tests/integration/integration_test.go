// +build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orchard9/pg-goer/internal/analyzer"
	"github.com/orchard9/pg-goer/internal/reporter"
	"github.com/orchard9/pg-goer/pkg/models"
)

const (
	testConnectionString = "postgresql://test_user:test_password@127.0.0.1:5556/integration_test?sslmode=disable"
	dockerComposeFile    = "docker-compose.yml"
	maxRetries           = 30
	retryInterval        = 2 * time.Second
)

func TestMain(m *testing.M) {
	// Setup: Start Docker Compose
	if err := startDockerCompose(); err != nil {
		panic("Failed to start Docker Compose: " + err.Error())
	}

	// Wait for PostgreSQL to be ready
	if err := waitForPostgreSQL(); err != nil {
		stopDockerCompose()
		panic("PostgreSQL failed to start: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup: Stop Docker Compose
	stopDockerCompose()

	os.Exit(code)
}

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	conn, err := analyzer.Connect(ctx, testConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer conn.Close()

	// Verify connection works by attempting a simple query
	if err := conn.DB().PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestSchemaAnalysis(t *testing.T) {
	ctx := context.Background()

	conn, err := analyzer.Connect(ctx, testConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer conn.Close()

	schemaAnalyzer := analyzer.NewSchemaAnalyzer(conn)

	// Test table discovery
	tables, err := schemaAnalyzer.GetTables(ctx, []string{"public"})
	if err != nil {
		t.Fatalf("Failed to get tables: %v", err)
	}

	expectedTables := map[string]bool{
		"users":          true,
		"posts":          true,
		"comments":       true,
		"categories":     true,
		"post_categories": true,
	}

	if len(tables) != len(expectedTables) {
		t.Errorf("Expected %d tables, got %d", len(expectedTables), len(tables))
	}

	foundTables := make(map[string]bool)
	for _, table := range tables {
		foundTables[table.Name] = true
	}

	for expected := range expectedTables {
		if !foundTables[expected] {
			t.Errorf("Expected table '%s' not found", expected)
		}
	}

	// Test column analysis for users table
	var usersTable *models.Table
	for i := range tables {
		if tables[i].Name == "users" {
			usersTable = &tables[i]
			break
		}
	}

	if usersTable == nil {
		t.Fatal("Users table not found")
	}

	columns, err := schemaAnalyzer.GetColumns(ctx, usersTable)
	if err != nil {
		t.Fatalf("Failed to get columns for users table: %v", err)
	}

	usersTable.Columns = columns

	// Verify key columns exist
	expectedColumns := map[string]bool{
		"id":         true,
		"username":   true,
		"email":      true,
		"full_name":  true,
		"age":        true,
		"is_active":  true,
		"created_at": true,
		"metadata":   true,
	}

	foundColumns := make(map[string]bool)
	for _, col := range columns {
		foundColumns[col.Name] = true
	}

	for expected := range expectedColumns {
		if !foundColumns[expected] {
			t.Errorf("Expected column '%s' not found in users table", expected)
		}
	}

	// Test primary key detection
	var idColumn *models.Column
	for i := range columns {
		if columns[i].Name == "id" {
			idColumn = &columns[i]
			break
		}
	}

	if idColumn == nil || !idColumn.IsPrimaryKey {
		t.Error("ID column should be marked as primary key")
	}

	// Test foreign key analysis
	foreignKeys, err := schemaAnalyzer.GetForeignKeys(ctx, usersTable)
	if err != nil {
		t.Fatalf("Failed to get foreign keys for users table: %v", err)
	}

	// Users table should have no foreign keys
	if len(foreignKeys) != 0 {
		t.Errorf("Users table should have no foreign keys, got %d", len(foreignKeys))
	}

	// Test foreign keys on posts table
	var postsTable *models.Table
	for i := range tables {
		if tables[i].Name == "posts" {
			postsTable = &tables[i]
			break
		}
	}

	if postsTable == nil {
		t.Fatal("Posts table not found")
	}

	postsForeignKeys, err := schemaAnalyzer.GetForeignKeys(ctx, postsTable)
	if err != nil {
		t.Fatalf("Failed to get foreign keys for posts table: %v", err)
	}

	// Posts table should have one foreign key to users
	if len(postsForeignKeys) != 1 {
		t.Errorf("Posts table should have 1 foreign key, got %d", len(postsForeignKeys))
	}

	if len(postsForeignKeys) > 0 {
		fk := postsForeignKeys[0]
		if fk.SourceColumn != "author_id" {
			t.Errorf("Expected foreign key on author_id, got %s", fk.SourceColumn)
		}
		if !strings.Contains(fk.ReferencedTable, "users") {
			t.Errorf("Expected foreign key to reference users table, got %s", fk.ReferencedTable)
		}
	}
}

func TestRowCountAnalysis(t *testing.T) {
	ctx := context.Background()

	conn, err := analyzer.Connect(ctx, testConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer conn.Close()

	schemaAnalyzer := analyzer.NewSchemaAnalyzer(conn)

	tables, err := schemaAnalyzer.GetTables(ctx, []string{"public"})
	if err != nil {
		t.Fatalf("Failed to get tables: %v", err)
	}

	rowCounts, err := schemaAnalyzer.GetTableRowCounts(ctx, tables)
	if err != nil {
		t.Fatalf("Failed to get row counts: %v", err)
	}

	// Verify we have row counts for expected tables
	expectedRowCounts := map[string]int64{
		"users":          3,
		"posts":          3,
		"comments":       4,
		"categories":     4,
		"post_categories": 4,
	}

	for tableName, expectedCount := range expectedRowCounts {
		if count, exists := rowCounts[tableName]; !exists {
			t.Errorf("Row count not found for table '%s'", tableName)
		} else if count != expectedCount {
			t.Errorf("Expected %d rows in '%s', got %d", expectedCount, tableName, count)
		}
	}
}

func TestCompleteWorkflow(t *testing.T) {
	ctx := context.Background()

	conn, err := analyzer.Connect(ctx, testConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer conn.Close()

	schemaAnalyzer := analyzer.NewSchemaAnalyzer(conn)

	// Get all schema information
	tables, err := schemaAnalyzer.GetTables(ctx, []string{"public"})
	if err != nil {
		t.Fatalf("Failed to get tables: %v", err)
	}

	// Populate tables with column and foreign key information
	for i := range tables {
		columns, err := schemaAnalyzer.GetColumns(ctx, &tables[i])
		if err != nil {
			t.Fatalf("Failed to get columns for %s: %v", tables[i].Name, err)
		}
		tables[i].Columns = columns

		foreignKeys, err := schemaAnalyzer.GetForeignKeys(ctx, &tables[i])
		if err != nil {
			t.Fatalf("Failed to get foreign keys for %s: %v", tables[i].Name, err)
		}
		tables[i].ForeignKeys = foreignKeys
	}

	// Get row counts
	rowCounts, err := schemaAnalyzer.GetTableRowCounts(ctx, tables)
	if err != nil {
		t.Fatalf("Failed to get row counts: %v", err)
	}

	// Apply row counts to tables
	for i := range tables {
		if count, exists := rowCounts[tables[i].Name]; exists {
			tables[i].RowCount = count
		}
	}

	// Create schema model
	schema := models.Schema{
		Name:       "Integration Test Database",
		Tables:     tables,
		Views:      []models.View{},
		Sequences:  []models.Sequence{},
		Extensions: []models.Extension{},
	}

	// Generate markdown documentation
	markdownReporter := reporter.NewMarkdownReporter()
	documentation, err := markdownReporter.Generate(&schema)
	if err != nil {
		t.Fatalf("Failed to generate documentation: %v", err)
	}

	// Verify documentation content
	requiredSections := []string{
		"# PostgreSQL Database Documentation",
		"## Table of Contents",
		"## Database Summary",
		"## Database Relationships", // Should be present due to foreign keys
		"```mermaid",
		"erDiagram",
		"## Tables",
		"## users",
		"## posts",
		"## comments",
		"## categories",
		"## post_categories",
	}

	for _, section := range requiredSections {
		if !strings.Contains(documentation, section) {
			t.Errorf("Documentation missing required section: %s", section)
		}
	}

	// Verify relationships are documented
	expectedRelationships := []string{
		"users ||--o{ posts",
		"posts ||--o{ comments",
		"users ||--o{ comments",
		"categories ||--o{ categories",
		"posts ||--o{ post_categories",
		"categories ||--o{ post_categories",
	}

	for _, relationship := range expectedRelationships {
		if !strings.Contains(documentation, relationship) {
			t.Errorf("Documentation missing relationship: %s", relationship)
		}
	}

	// Verify row counts are included
	if !strings.Contains(documentation, "Row Count:") {
		t.Error("Documentation should include row counts")
	}

	// Verify the documentation is substantial
	if len(documentation) < 2000 {
		t.Errorf("Documentation seems too short: %d characters", len(documentation))
	}
}

// Helper functions

func startDockerCompose() error {
	cmd := exec.Command("docker", "compose", "-f", dockerComposeFile, "up", "-d")
	cmd.Dir = getIntegrationDir()
	return cmd.Run()
}

func stopDockerCompose() {
	cmd := exec.Command("docker", "compose", "-f", dockerComposeFile, "down", "-v", "--remove-orphans")
	cmd.Dir = getIntegrationDir()
	cmd.Run() // Ignore errors during cleanup
}

func waitForPostgreSQL() error {
	for i := 0; i < maxRetries; i++ {
		ctx := context.Background()
		conn, err := analyzer.Connect(ctx, testConnectionString)
		if err == nil {
			if err := conn.DB().PingContext(ctx); err == nil {
				conn.Close()
				return nil
			}
			conn.Close()
		}
		time.Sleep(retryInterval)
	}
	return nil
}

func getIntegrationDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// If we're already in the integration directory, return current dir
	if filepath.Base(wd) == "integration" {
		return "."
	}

	// Look for integration directory
	integrationDir := filepath.Join(wd, "tests", "integration")
	if _, err := os.Stat(integrationDir); err == nil {
		return integrationDir
	}

	// Try relative path from project root
	integrationDir = filepath.Join("..", "..", "tests", "integration")
	if _, err := os.Stat(integrationDir); err == nil {
		return integrationDir
	}

	// Default to current directory
	return "."
}