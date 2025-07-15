package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/orchard9/pg-goer/internal/analyzer"
	"github.com/orchard9/pg-goer/internal/reporter"
	"github.com/orchard9/pg-goer/pkg/models"
)

const (
	defaultOutput    = "database-docs.md"
	defaultFormat    = "markdown"
	defaultTimeout   = 10 * time.Second
	defaultMaxTables = 1000
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		output     string
		format     string
		schemas    string
		showHelp   bool
		verbose    bool
		versionCmd bool
	)

	flag.StringVar(&output, "output", defaultOutput, "Output file")
	flag.StringVar(&output, "o", defaultOutput, "Output file (shorthand)")
	flag.StringVar(&format, "format", defaultFormat, "Output format (markdown or json)")
	flag.StringVar(&format, "f", defaultFormat, "Output format (shorthand)")
	flag.StringVar(&schemas, "schemas", "", "Comma-separated list of schemas to document")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.BoolVar(&versionCmd, "version", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "PG Go ER - PostgreSQL database documentation generator\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  pg-goer [flags] <connection-string>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  pg-goer \"postgresql://user:password@localhost/dbname\"\n")
		fmt.Fprintf(os.Stderr, "  pg-goer -schemas public,custom \"postgresql://localhost/mydb\"\n")
		fmt.Fprintf(os.Stderr, "  pg-goer -format json -output db.json \"postgresql://localhost/mydb\"\n")
	}

	flag.Parse()

	if versionCmd {
		fmt.Printf("pg-goer version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	connectionString := ""
	if flag.NArg() > 0 {
		connectionString = flag.Arg(0)
	} else if env := os.Getenv("DATABASE_URL"); env != "" {
		connectionString = env
	} else {
		fmt.Fprintf(os.Stderr, "Error: connection string required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var schemaList []string
	if schemas != "" {
		schemaList = strings.Split(schemas, ",")
		for i := range schemaList {
			schemaList[i] = strings.TrimSpace(schemaList[i])
		}
	}

	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(os.Stderr)
	}

	if err := run(connectionString, output, format, schemaList); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(connectionString, output, format string, schemas []string) error {
	// Validate format
	if format != "markdown" && format != "json" {
		return fmt.Errorf("invalid format '%s': must be 'markdown' or 'json'", format)
	}

	ctx := context.Background()

	conn, schemaAnalyzer, err := connectToDatabase(ctx, connectionString)
	if err != nil {
		return err
	}

	defer conn.Close()

	tables, err := fetchAllTableData(ctx, schemaAnalyzer, schemas)
	if err != nil {
		return err
	}

	extensions, err := fetchExtensions(ctx, schemaAnalyzer)
	if err != nil {
		return err
	}

	views, err := fetchViews(ctx, schemaAnalyzer, schemas)
	if err != nil {
		return err
	}

	sequences, err := fetchSequences(ctx, schemaAnalyzer, schemas)
	if err != nil {
		return err
	}

	schema := models.Schema{
		Name:       "Database Documentation",
		Tables:     tables,
		Views:      views,
		Sequences:  sequences,
		Extensions: extensions,
	}

	return generateAndWriteDocumentation(&schema, format, output)
}

func connectToDatabase(ctx context.Context, connectionString string) (*analyzer.Connection, *analyzer.SchemaAnalyzer, error) {
	log.Println("Connecting to database...")

	conn, err := analyzer.Connect(ctx, connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, analyzer.NewSchemaAnalyzer(conn), nil
}

func fetchAllTableData(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, schemas []string) ([]models.Table, error) {
	log.Println("Fetching tables...")

	tables, err := schemaAnalyzer.GetTables(ctx, schemas)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}

	log.Printf("Found %d tables\n", len(tables))

	if err := enrichTablesWithMetadata(ctx, schemaAnalyzer, tables); err != nil {
		return nil, err
	}

	if err := addRowCounts(ctx, schemaAnalyzer, tables); err != nil {
		return nil, err
	}

	return tables, nil
}

func fetchExtensions(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer) ([]models.Extension, error) {
	log.Println("Fetching PostgreSQL extensions...")

	extensions, err := schemaAnalyzer.GetExtensions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get extensions: %w", err)
	}

	return extensions, nil
}

func fetchViews(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, schemas []string) ([]models.View, error) {
	log.Println("Fetching views...")

	views, err := schemaAnalyzer.GetViews(ctx, schemas)
	if err != nil {
		return nil, fmt.Errorf("failed to get views: %w", err)
	}

	return views, nil
}

func fetchSequences(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, schemas []string) ([]models.Sequence, error) {
	log.Println("Fetching sequences...")

	sequences, err := schemaAnalyzer.GetSequences(ctx, schemas)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequences: %w", err)
	}

	return sequences, nil
}

func enrichTablesWithMetadata(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, tables []models.Table) error {
	for i := range tables {
		if err := fetchTableColumns(ctx, schemaAnalyzer, &tables[i]); err != nil {
			return err
		}

		if err := fetchTableForeignKeys(ctx, schemaAnalyzer, &tables[i]); err != nil {
			return err
		}

		if err := fetchTableIndexes(ctx, schemaAnalyzer, &tables[i]); err != nil {
			return err
		}

		if err := fetchTableTriggers(ctx, schemaAnalyzer, &tables[i]); err != nil {
			return err
		}
	}

	return nil
}

func fetchTableColumns(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, table *models.Table) error {
	log.Printf("Fetching columns for %s.%s...\n", table.Schema, table.Name)

	columns, err := schemaAnalyzer.GetColumns(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to get columns for %s.%s: %w", table.Schema, table.Name, err)
	}

	table.Columns = columns

	return nil
}

func fetchTableForeignKeys(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, table *models.Table) error {
	log.Printf("Fetching foreign keys for %s.%s...\n", table.Schema, table.Name)

	foreignKeys, err := schemaAnalyzer.GetForeignKeys(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to get foreign keys for %s.%s: %w", table.Schema, table.Name, err)
	}

	table.ForeignKeys = foreignKeys

	return nil
}

func fetchTableIndexes(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, table *models.Table) error {
	log.Printf("Fetching indexes for %s.%s...\n", table.Schema, table.Name)

	indexes, err := schemaAnalyzer.GetIndexes(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to get indexes for %s.%s: %w", table.Schema, table.Name, err)
	}

	table.Indexes = indexes

	return nil
}

func fetchTableTriggers(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, table *models.Table) error {
	log.Printf("Fetching triggers for %s.%s...\n", table.Schema, table.Name)

	triggers, err := schemaAnalyzer.GetTriggers(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to get triggers for %s.%s: %w", table.Schema, table.Name, err)
	}

	table.Triggers = triggers

	return nil
}

func addRowCounts(ctx context.Context, schemaAnalyzer *analyzer.SchemaAnalyzer, tables []models.Table) error {
	log.Println("Fetching table row counts...")

	rowCounts, err := schemaAnalyzer.GetTableRowCounts(ctx, tables)
	if err != nil {
		return fmt.Errorf("failed to get table row counts: %w", err)
	}

	// Apply row counts to tables
	for i := range tables {
		if count, exists := rowCounts[tables[i].Name]; exists {
			tables[i].RowCount = count
		}
	}

	return nil
}

func generateAndWriteDocumentation(schema *models.Schema, format, output string) error {
	log.Println("Generating documentation...")

	documentation, err := generateDocumentation(schema, format)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, []byte(documentation), 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Printf("Documentation written to %s\n", output)

	return nil
}

func generateDocumentation(schema *models.Schema, format string) (string, error) {
	switch format {
	case "markdown":
		markdownReporter := reporter.NewMarkdownReporter()
		return markdownReporter.Generate(schema)
	case "json":
		jsonReporter := reporter.NewJSONReporter()
		return jsonReporter.Generate(schema)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
