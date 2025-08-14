package analyzer

import (
	"fmt"
	"strings"
)

// DatabaseType represents the type of database being analyzed.
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgresql"
	MariaDB    DatabaseType = "mariadb"
)

// DetectDatabaseType detects the database type from a connection string.
func DetectDatabaseType(connectionString string) (DatabaseType, error) {
	lower := strings.ToLower(connectionString)

	// Check for PostgreSQL patterns
	if strings.HasPrefix(lower, "postgresql://") ||
		strings.HasPrefix(lower, "postgres://") ||
		strings.Contains(lower, "host=") && strings.Contains(lower, "dbname=") {
		return PostgreSQL, nil
	}

	// Check for MariaDB/MySQL patterns
	if strings.HasPrefix(lower, "mysql://") ||
		strings.HasPrefix(lower, "mariadb://") ||
		strings.Contains(lower, "@tcp(") {
		return MariaDB, nil
	}

	return "", fmt.Errorf("unable to detect database type from connection string")
}

// ParseConnectionString returns the driver name and connection string for the database.
func ParseConnectionString(dbType DatabaseType, connectionString string) (driver string, connStr string, err error) {
	switch dbType {
	case PostgreSQL:
		return "postgres", connectionString, nil
	case MariaDB:
		// Handle both URL and DSN formats for MySQL/MariaDB
		if strings.HasPrefix(strings.ToLower(connectionString), "mysql://") ||
			strings.HasPrefix(strings.ToLower(connectionString), "mariadb://") {
			// Convert URL format to DSN format
			connStr, err = convertMySQLURLToDSN(connectionString)
			if err != nil {
				return "", "", err
			}
			return "mysql", connStr, nil
		}
		// Already in DSN format
		return "mysql", connectionString, nil
	default:
		return "", "", fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// convertMySQLURLToDSN converts a MySQL URL to DSN format
// mysql://user:password@host:port/database -> user:password@tcp(host:port)/database
func convertMySQLURLToDSN(url string) (string, error) {
	// Remove protocol prefix
	url = strings.TrimPrefix(url, "mysql://")
	url = strings.TrimPrefix(url, "mariadb://")

	// Parse components
	parts := strings.SplitN(url, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid MySQL URL format")
	}

	authAndHost := parts[0]
	database := parts[1]

	// Split auth and host
	authParts := strings.SplitN(authAndHost, "@", 2)
	if len(authParts) != 2 {
		return "", fmt.Errorf("invalid MySQL URL format: missing @ separator")
	}

	auth := authParts[0]
	host := authParts[1]

	// Add default port if not specified
	if !strings.Contains(host, ":") {
		host += ":3306"
	}

	// Build DSN
	dsn := fmt.Sprintf("%s@tcp(%s)/%s?parseTime=true", auth, host, database)

	return dsn, nil
}
