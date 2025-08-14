package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/orchard9/pg-goer/pkg/models"
)

type MariaDBAnalyzer struct {
	conn *Connection
}

func (a *MariaDBAnalyzer) GetTables(ctx context.Context, schemas []string) ([]models.Table, error) {
	query := a.buildTableQuery(schemas)

	return querySchemaObjects(ctx, a.conn.db, query, schemas, func() models.Table { return models.Table{} },
		func(item *models.Table) []interface{} { return []interface{}{&item.Schema, &item.Name} },
		"tables")
}

func (a *MariaDBAnalyzer) buildTableQuery(schemas []string) string {
	return a.buildSchemaFilterQuery(
		`SELECT table_schema AS schema_name, table_name 
		 FROM information_schema.tables 
		 WHERE table_type = 'BASE TABLE'`,
		"table_schema",
		"table_schema, table_name",
		schemas,
	)
}

func (a *MariaDBAnalyzer) GetColumns(ctx context.Context, table *models.Table) ([]models.Column, error) {
	query := `
		SELECT 
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			c.character_maximum_length,
			CASE WHEN c.column_key = 'PRI' THEN true ELSE false END AS is_primary_key,
			CASE WHEN c.column_key IN ('UNI', 'PRI') THEN true ELSE false END AS is_unique
		FROM 
			information_schema.columns c
		WHERE 
			c.table_schema = ?
			AND c.table_name = ?
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

func (a *MariaDBAnalyzer) GetForeignKeys(ctx context.Context, table *models.Table) ([]models.ForeignKey, error) {
	query := `
		SELECT 
			tc.constraint_name,
			kcu.column_name,
			kcu.referenced_table_schema,
			kcu.referenced_table_name,
			kcu.referenced_column_name,
			rc.delete_rule,
			rc.update_rule
		FROM 
			information_schema.table_constraints AS tc 
		JOIN 
			information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN 
			information_schema.referential_constraints AS rc
			ON tc.constraint_name = rc.constraint_name
			AND tc.table_schema = rc.constraint_schema
		WHERE 
			tc.constraint_type = 'FOREIGN KEY' 
			AND tc.table_schema = ?
			AND tc.table_name = ?
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

func (a *MariaDBAnalyzer) GetTableRowCounts(ctx context.Context, tables []models.Table) (map[string]int64, error) {
	if len(tables) == 0 {
		return make(map[string]int64), nil
	}

	rowCounts := make(map[string]int64)

	// For MariaDB, we'll use table_rows from information_schema.tables
	// Note: This is an approximation for InnoDB tables
	for _, table := range tables {
		query := `
			SELECT COALESCE(table_rows, 0) 
			FROM information_schema.tables 
			WHERE table_schema = ? AND table_name = ?`

		var rowCount int64
		err := a.conn.db.QueryRowContext(ctx, query, table.Schema, table.Name).Scan(&rowCount)
		if err != nil {
			// If we can't get the estimate, try COUNT(*) as fallback
			// Note: This is safe as table names are validated by the database schema
			countQuery := "SELECT COUNT(*) FROM `" + table.Schema + "`.`" + table.Name + "`"
			err = a.conn.db.QueryRowContext(ctx, countQuery).Scan(&rowCount)
			if err != nil {
				// If both fail, set to 0
				rowCount = 0
			}
		}

		rowCounts[table.Name] = rowCount
	}

	return rowCounts, nil
}

func (a *MariaDBAnalyzer) GetIndexes(ctx context.Context, table *models.Table) ([]models.Index, error) {
	query := `
		SELECT DISTINCT
			s.index_name,
			CASE 
				WHEN s.index_name = 'PRIMARY' THEN 'PRIMARY KEY'
				WHEN s.non_unique = 0 THEN 'UNIQUE'
				ELSE 'INDEX'
			END AS index_type,
			CASE WHEN s.index_name = 'PRIMARY' THEN true ELSE false END AS is_primary,
			CASE WHEN s.non_unique = 0 THEN true ELSE false END AS is_unique,
			s.index_type AS access_method,
			GROUP_CONCAT(s.column_name ORDER BY s.seq_in_index SEPARATOR ',') AS columns
		FROM 
			information_schema.statistics s
		WHERE 
			s.table_schema = ?
			AND s.table_name = ?
		GROUP BY 
			s.index_name, s.non_unique, s.index_type
		ORDER BY 
			is_primary DESC, is_unique DESC, s.index_name`

	rows, err := a.conn.db.QueryContext(ctx, query, table.Schema, table.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var indexes []models.Index

	for rows.Next() {
		var (
			idx     models.Index
			columns string
		)

		if err := rows.Scan(
			&idx.Name,
			&idx.Type,
			&idx.IsPrimary,
			&idx.IsUnique,
			&idx.Method,
			&columns,
		); err != nil {
			return nil, fmt.Errorf("failed to scan index row: %w", err)
		}

		if columns != "" {
			idx.Columns = strings.Split(columns, ",")
		}

		indexes = append(indexes, idx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating index rows: %w", err)
	}

	return indexes, nil
}

func (a *MariaDBAnalyzer) GetTriggers(ctx context.Context, table *models.Table) ([]models.Trigger, error) {
	query := `
		SELECT 
			t.trigger_name,
			t.action_timing AS timing,
			t.event_manipulation AS event,
			t.action_statement AS function_name,
			'ROW' AS orientation  -- MariaDB triggers are always row-level
		FROM 
			information_schema.triggers t
		WHERE 
			t.trigger_schema = ?
			AND t.event_object_table = ?
		ORDER BY 
			t.trigger_name`

	rows, err := a.conn.db.QueryContext(ctx, query, table.Schema, table.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to query triggers: %w", err)
	}
	defer rows.Close()

	var triggers []models.Trigger

	for rows.Next() {
		var trigger models.Trigger

		if err := rows.Scan(
			&trigger.Name,
			&trigger.Timing,
			&trigger.Event,
			&trigger.Function,
			&trigger.Orientation,
		); err != nil {
			return nil, fmt.Errorf("failed to scan trigger row: %w", err)
		}

		triggers = append(triggers, trigger)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trigger rows: %w", err)
	}

	return triggers, nil
}

func (a *MariaDBAnalyzer) GetExtensions(_ context.Context) ([]models.Extension, error) {
	// MariaDB doesn't have extensions like PostgreSQL
	// Return empty slice
	return []models.Extension{}, nil
}

func (a *MariaDBAnalyzer) GetViews(ctx context.Context, schemas []string) ([]models.View, error) {
	query := a.buildViewQuery(schemas)

	return querySchemaObjects(ctx, a.conn.db, query, schemas, func() models.View { return models.View{} },
		func(item *models.View) []interface{} { return []interface{}{&item.Schema, &item.Name} },
		"views")
}

func (a *MariaDBAnalyzer) buildViewQuery(schemas []string) string {
	return a.buildSchemaFilterQuery(
		`SELECT table_schema AS schema_name, table_name AS view_name 
		 FROM information_schema.views WHERE true`,
		"table_schema",
		"table_schema, table_name",
		schemas,
	)
}

func (a *MariaDBAnalyzer) GetSequences(_ context.Context, _ []string) ([]models.Sequence, error) {
	// MariaDB doesn't have sequences like PostgreSQL
	// Return empty slice
	return []models.Sequence{}, nil
}

func (a *MariaDBAnalyzer) buildSchemaFilterQuery(baseQuery, schemaColumn, orderBy string, schemas []string) string {
	whereClause := " AND "

	switch len(schemas) {
	case 0:
		whereClause += fmt.Sprintf("%s NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')", schemaColumn)
	case 1:
		whereClause += fmt.Sprintf("%s = ?", schemaColumn)
	default:
		placeholders := make([]string, len(schemas))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		whereClause += fmt.Sprintf("%s IN (%s)", schemaColumn, strings.Join(placeholders, ", "))
	}

	return baseQuery + whereClause + fmt.Sprintf(" ORDER BY %s", orderBy)
}
