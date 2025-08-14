package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL/MariaDB driver
	_ "github.com/lib/pq"              // PostgreSQL driver
)

type Connection struct {
	db     *sql.DB
	dbType DatabaseType
}

func Connect(ctx context.Context, connectionString string) (*Connection, error) {
	// Auto-detect database type
	dbType, err := DetectDatabaseType(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to detect database type: %w", err)
	}

	return ConnectWithType(ctx, connectionString, dbType)
}

func ConnectWithType(ctx context.Context, connectionString string, dbType DatabaseType) (*Connection, error) {
	driver, connStr, err := ParseConnectionString(dbType, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{db: db, dbType: dbType}, nil
}

func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) DatabaseType() DatabaseType {
	return c.dbType
}
