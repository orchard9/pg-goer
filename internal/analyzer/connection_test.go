package analyzer

import (
	"context"
	"testing"
)

func TestConnectionStringParsing(t *testing.T) {
	tests := []struct {
		name             string
		connectionString string
		shouldFail       bool
	}{
		{
			name:             "valid URI format",
			connectionString: "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			shouldFail:       false,
		},
		{
			name:             "valid key-value format",
			connectionString: "host=localhost port=5432 dbname=testdb user=postgres password=secret sslmode=disable",
			shouldFail:       false,
		},
		{
			name:             "URI with special characters in password",
			connectionString: "postgres://user:p@ss%40word@localhost:5432/dbname",
			shouldFail:       false,
		},
		{
			name:             "minimal connection string",
			connectionString: "dbname=testdb",
			shouldFail:       false,
		},
		{
			name:             "empty connection string",
			connectionString: "",
			shouldFail:       true,
		},
		{
			name:             "invalid format",
			connectionString: "not a valid connection string",
			shouldFail:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := Connect(context.Background(), tt.connectionString)

			if tt.shouldFail && err == nil {
				t.Errorf("expected connection to fail but succeeded")
			}

			if conn != nil {
				conn.Close()
			}
		})
	}
}
