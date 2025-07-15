package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/orchard9/pg-goer/pkg/models"
)

type SchemaAnalyzer struct {
	conn *Connection
}

func NewSchemaAnalyzer(conn *Connection) *SchemaAnalyzer {
	return &SchemaAnalyzer{conn: conn}
}

func (a *SchemaAnalyzer) GetTables(ctx context.Context, schemas []string) ([]models.Table, error) {
	query := a.buildTableQuery(schemas)
	args := make([]interface{}, len(schemas))

	for i, schema := range schemas {
		args[i] = schema
	}

	rows, err := a.conn.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []models.Table

	for rows.Next() {
		var table models.Table
		if err := rows.Scan(&table.Schema, &table.Name); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}

		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table rows: %w", err)
	}

	return tables, nil
}

func (a *SchemaAnalyzer) buildTableQuery(schemas []string) string {
	baseQuery := `
		SELECT 
			n.nspname AS schema_name,
			c.relname AS table_name
		FROM 
			pg_catalog.pg_class c
		INNER JOIN 
			pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE 
			c.relkind = 'r'
			AND `

	if len(schemas) == 0 {
		return baseQuery + `NOT n.nspname IN ('pg_catalog', 'information_schema', 'pg_toast')
		ORDER BY n.nspname, c.relname`
	}

	if len(schemas) == 1 {
		return baseQuery + `n.nspname = $1
		ORDER BY n.nspname, c.relname`
	}

	placeholders := make([]string, len(schemas))
	for i := range schemas {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return baseQuery + fmt.Sprintf(`n.nspname IN (%s)
		ORDER BY n.nspname, c.relname`, strings.Join(placeholders, ", "))
}

func (a *SchemaAnalyzer) GetColumns(ctx context.Context, table *models.Table) ([]models.Column, error) {
	query := `
		SELECT 
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			c.character_maximum_length,
			COALESCE(
				(SELECT true 
				 FROM information_schema.table_constraints tc
				 JOIN information_schema.key_column_usage kcu 
				   ON tc.constraint_name = kcu.constraint_name 
				   AND tc.table_schema = kcu.table_schema
				 WHERE tc.constraint_type = 'PRIMARY KEY' 
				   AND tc.table_schema = c.table_schema 
				   AND tc.table_name = c.table_name 
				   AND kcu.column_name = c.column_name
				), false
			) AS is_primary_key,
			COALESCE(
				(SELECT true 
				 FROM information_schema.table_constraints tc
				 JOIN information_schema.key_column_usage kcu 
				   ON tc.constraint_name = kcu.constraint_name 
				   AND tc.table_schema = kcu.table_schema
				 WHERE tc.constraint_type = 'UNIQUE' 
				   AND tc.table_schema = c.table_schema 
				   AND tc.table_name = c.table_name 
				   AND kcu.column_name = c.column_name
				), false
			) AS is_unique
		FROM 
			information_schema.columns c
		WHERE 
			c.table_schema = $1
			AND c.table_name = $2
		ORDER BY 
			c.ordinal_position`

	rows, err := a.conn.db.QueryContext(ctx, query, table.Schema, table.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []models.Column

	for rows.Next() {
		var (
			col          models.Column
			isNullable   string
			defaultValue sql.NullString
			maxLength    sql.NullInt64
		)

		if err := rows.Scan(
			&col.Name,
			&col.DataType,
			&isNullable,
			&defaultValue,
			&maxLength,
			&col.IsPrimaryKey,
			&col.IsUnique,
		); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}

		col.IsNullable = isNullable == "YES"

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
		}

		if maxLength.Valid {
			length := int(maxLength.Int64)
			col.MaxLength = &length
		}

		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column rows: %w", err)
	}

	return columns, nil
}

func (a *SchemaAnalyzer) GetForeignKeys(ctx context.Context, table *models.Table) ([]models.ForeignKey, error) {
	query := `
		SELECT 
			tc.constraint_name,
			kcu.column_name,
			ccu.table_schema AS foreign_table_schema,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.delete_rule,
			rc.update_rule
		FROM 
			information_schema.table_constraints AS tc 
		JOIN 
			information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN 
			information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		JOIN 
			information_schema.referential_constraints AS rc
			ON tc.constraint_name = rc.constraint_name
			AND tc.table_schema = rc.constraint_schema
		WHERE 
			tc.constraint_type = 'FOREIGN KEY' 
			AND tc.table_schema = $1
			AND tc.table_name = $2
		ORDER BY 
			tc.constraint_name, kcu.ordinal_position`

	rows, err := a.conn.db.QueryContext(ctx, query, table.Schema, table.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign keys: %w", err)
	}
	defer rows.Close()

	var foreignKeys []models.ForeignKey

	for rows.Next() {
		var (
			fk            models.ForeignKey
			foreignSchema string
		)

		if err := rows.Scan(
			&fk.Name,
			&fk.SourceColumn,
			&foreignSchema,
			&fk.ReferencedTable,
			&fk.ReferencedColumn,
			&fk.OnDelete,
			&fk.OnUpdate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan foreign key row: %w", err)
		}

		fk.SourceTable = table.Name
		fk.ReferencedTable = fmt.Sprintf("%s.%s", foreignSchema, fk.ReferencedTable)

		foreignKeys = append(foreignKeys, fk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foreign key rows: %w", err)
	}

	return foreignKeys, nil
}

func (a *SchemaAnalyzer) GetTableRowCounts(ctx context.Context, tables []models.Table) (map[string]int64, error) {
	if len(tables) == 0 {
		return make(map[string]int64), nil
	}

	// Build table name list for query
	tableNames := make([]string, 0, len(tables))
	tableMap := make(map[string]string) // qualified name -> simple name

	for _, table := range tables {
		qualifiedName := fmt.Sprintf("%s.%s", table.Schema, table.Name)
		tableNames = append(tableNames, qualifiedName)
		tableMap[qualifiedName] = table.Name
	}

	query := `
		SELECT 
			schemaname || '.' || relname AS qualified_name,
			COALESCE(n_tup_ins - n_tup_del, 0) AS row_count
		FROM 
			pg_stat_user_tables 
		WHERE 
			schemaname || '.' || relname = ANY($1)`

	rows, err := a.conn.db.QueryContext(ctx, query, "{"+strings.Join(tableNames, ",")+"}")
	if err != nil {
		return nil, fmt.Errorf("failed to query table row counts: %w", err)
	}
	defer rows.Close()

	rowCounts := make(map[string]int64)

	for rows.Next() {
		var (
			qualifiedName string
			rowCount      int64
		)

		if err := rows.Scan(&qualifiedName, &rowCount); err != nil {
			return nil, fmt.Errorf("failed to scan row count row: %w", err)
		}

		if simpleName, exists := tableMap[qualifiedName]; exists {
			rowCounts[simpleName] = rowCount
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating row count rows: %w", err)
	}

	return rowCounts, nil
}
