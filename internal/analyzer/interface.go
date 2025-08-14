package analyzer

import (
	"context"
	"fmt"

	"github.com/orchard9/pg-goer/pkg/models"
)

// DatabaseAnalyzer defines the interface for analyzing database schemas
type DatabaseAnalyzer interface {
	// GetTables returns all tables in the specified schemas
	GetTables(ctx context.Context, schemas []string) ([]models.Table, error)

	// GetColumns returns all columns for a specific table
	GetColumns(ctx context.Context, table *models.Table) ([]models.Column, error)

	// GetForeignKeys returns all foreign keys for a specific table
	GetForeignKeys(ctx context.Context, table *models.Table) ([]models.ForeignKey, error)

	// GetIndexes returns all indexes for a specific table
	GetIndexes(ctx context.Context, table *models.Table) ([]models.Index, error)

	// GetTriggers returns all triggers for a specific table
	GetTriggers(ctx context.Context, table *models.Table) ([]models.Trigger, error)

	// GetTableRowCounts returns row counts for the specified tables
	GetTableRowCounts(ctx context.Context, tables []models.Table) (map[string]int64, error)

	// GetExtensions returns all database extensions (PostgreSQL specific)
	GetExtensions(ctx context.Context) ([]models.Extension, error)

	// GetViews returns all views in the specified schemas
	GetViews(ctx context.Context, schemas []string) ([]models.View, error)

	// GetSequences returns all sequences in the specified schemas
	GetSequences(ctx context.Context, schemas []string) ([]models.Sequence, error)
}

// NewDatabaseAnalyzer creates a new analyzer based on the database type
func NewDatabaseAnalyzer(dbType DatabaseType, conn *Connection) (DatabaseAnalyzer, error) {
	switch dbType {
	case PostgreSQL:
		return &PostgreSQLAnalyzer{conn: conn}, nil
	case MariaDB:
		return &MariaDBAnalyzer{conn: conn}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
