package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/orchard9/pg-goer/pkg/models"
)

type PostgreSQLAnalyzer struct {
	conn *Connection
}

// SchemaAnalyzer is a compatibility alias for PostgreSQLAnalyzer
type SchemaAnalyzer = PostgreSQLAnalyzer

func NewSchemaAnalyzer(conn *Connection) *SchemaAnalyzer {
	return &PostgreSQLAnalyzer{conn: conn}
}

func (a *PostgreSQLAnalyzer) GetTables(ctx context.Context, schemas []string) ([]models.Table, error) {
	query := a.buildTableQuery(schemas)

	return querySchemaObjects(ctx, a.conn.db, query, schemas, func() models.Table { return models.Table{} },
		func(item *models.Table) []interface{} { return []interface{}{&item.Schema, &item.Name} },
		"tables")
}

func (a *PostgreSQLAnalyzer) buildTableQuery(schemas []string) string {
	return a.buildSchemaFilterQuery(
		`SELECT n.nspname AS schema_name, c.relname AS table_name 
		 FROM pg_catalog.pg_class c 
		 INNER JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace 
		 WHERE c.relkind = 'r'`,
		"n.nspname",
		"n.nspname, c.relname",
		schemas,
	)
}

func (a *PostgreSQLAnalyzer) GetColumns(ctx context.Context, table *models.Table) ([]models.Column, error) {
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
			EXISTS (
				SELECT 1
				FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu 
				  ON tc.constraint_name = kcu.constraint_name 
				  AND tc.table_schema = kcu.table_schema
				WHERE tc.constraint_type = 'UNIQUE' 
				  AND tc.table_schema = c.table_schema 
				  AND tc.table_name = c.table_name 
				  AND kcu.column_name = c.column_name
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

func (a *PostgreSQLAnalyzer) GetForeignKeys(ctx context.Context, table *models.Table) ([]models.ForeignKey, error) {
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

func (a *PostgreSQLAnalyzer) GetTableRowCounts(ctx context.Context, tables []models.Table) (map[string]int64, error) {
	if len(tables) == 0 {
		return make(map[string]int64), nil
	}

	// Build table name list for query
	tableNames := make([]string, 0, len(tables))
	tableMap := make(map[string]string) // qualified name -> simple name

	for i := range tables {
		table := &tables[i]
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

func (a *PostgreSQLAnalyzer) GetIndexes(ctx context.Context, table *models.Table) ([]models.Index, error) {
	query := `
		SELECT DISTINCT
			i.relname AS index_name,
			CASE 
				WHEN ic.indisprimary THEN 'PRIMARY KEY'
				WHEN ic.indisunique THEN 'UNIQUE'
				ELSE 'INDEX'
			END AS index_type,
			ic.indisprimary AS is_primary,
			ic.indisunique AS is_unique,
			am.amname AS access_method,
			array_agg(a.attname ORDER BY ic.indkey) AS columns
		FROM 
			pg_catalog.pg_index ic
		JOIN 
			pg_catalog.pg_class i ON i.oid = ic.indexrelid
		JOIN 
			pg_catalog.pg_class t ON t.oid = ic.indrelid
		JOIN 
			pg_catalog.pg_namespace n ON n.oid = t.relnamespace
		JOIN 
			pg_catalog.pg_am am ON am.oid = i.relam
		JOIN 
			pg_catalog.pg_attribute a ON a.attrelid = t.oid 
				AND a.attnum = ANY(ic.indkey)
		WHERE 
			n.nspname = $1
			AND t.relname = $2
			AND t.relkind = 'r'
		GROUP BY 
			i.relname, ic.indisprimary, ic.indisunique, am.amname
		ORDER BY 
			ic.indisprimary DESC, ic.indisunique DESC, i.relname`

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

		// Parse the PostgreSQL array format {col1,col2}
		columns = strings.Trim(columns, "{}")
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

func (a *PostgreSQLAnalyzer) GetTriggers(ctx context.Context, table *models.Table) ([]models.Trigger, error) {
	query := `
		SELECT 
			t.tgname AS trigger_name,
			CASE t.tgtype & cast(2 as int2)
				WHEN 0 THEN 'AFTER'
				ELSE 'BEFORE'
			END AS timing,
			CASE t.tgtype & cast(28 as int2)
				WHEN 4 THEN 'INSERT'
				WHEN 8 THEN 'DELETE'
				WHEN 16 THEN 'UPDATE'
				WHEN 12 THEN 'INSERT,DELETE'
				WHEN 20 THEN 'INSERT,UPDATE'
				WHEN 24 THEN 'DELETE,UPDATE'
				WHEN 28 THEN 'INSERT,DELETE,UPDATE'
			END AS event,
			p.proname AS function_name,
			CASE t.tgtype & cast(1 as int2)
				WHEN 0 THEN 'STATEMENT'
				ELSE 'ROW'
			END AS orientation
		FROM 
			pg_catalog.pg_trigger t
		JOIN 
			pg_catalog.pg_class c ON c.oid = t.tgrelid
		JOIN 
			pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		JOIN 
			pg_catalog.pg_proc p ON p.oid = t.tgfoid
		WHERE 
			n.nspname = $1
			AND c.relname = $2
			AND NOT t.tgisinternal
		ORDER BY 
			t.tgname`

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

func (a *PostgreSQLAnalyzer) GetExtensions(ctx context.Context) ([]models.Extension, error) {
	query := `
		SELECT 
			e.extname AS extension_name,
			e.extversion AS extension_version,
			n.nspname AS schema_name
		FROM 
			pg_catalog.pg_extension e
		JOIN 
			pg_catalog.pg_namespace n ON n.oid = e.extnamespace
		WHERE 
			e.extname NOT IN ('plpgsql')  -- Exclude built-in extensions
		ORDER BY 
			e.extname`

	rows, err := a.conn.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query extensions: %w", err)
	}
	defer rows.Close()

	var extensions []models.Extension

	for rows.Next() {
		var ext models.Extension

		if err := rows.Scan(&ext.Name, &ext.Version, &ext.Schema); err != nil {
			return nil, fmt.Errorf("failed to scan extension row: %w", err)
		}

		extensions = append(extensions, ext)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating extension rows: %w", err)
	}

	return extensions, nil
}

func (a *PostgreSQLAnalyzer) GetViews(ctx context.Context, schemas []string) ([]models.View, error) {
	query := a.buildViewQuery(schemas)

	return querySchemaObjects(ctx, a.conn.db, query, schemas, func() models.View { return models.View{} },
		func(item *models.View) []interface{} { return []interface{}{&item.Schema, &item.Name} },
		"views")
}

func (a *PostgreSQLAnalyzer) buildViewQuery(schemas []string) string {
	return a.buildSchemaFilterQuery(
		`SELECT schemaname AS schema_name, viewname AS view_name FROM pg_catalog.pg_views WHERE true`,
		"schemaname",
		"schemaname, viewname",
		schemas,
	)
}

func (a *PostgreSQLAnalyzer) GetSequences(ctx context.Context, schemas []string) ([]models.Sequence, error) {
	query := a.buildSequenceQuery(schemas)
	args := make([]interface{}, len(schemas))

	for i, schema := range schemas {
		args[i] = schema
	}

	rows, err := a.conn.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sequences: %w", err)
	}
	defer rows.Close()

	var sequences []models.Sequence

	for rows.Next() {
		var seq models.Sequence
		if err := rows.Scan(&seq.Schema, &seq.Name, &seq.DataType, &seq.StartValue, &seq.MinValue, &seq.MaxValue, &seq.Increment); err != nil {
			return nil, fmt.Errorf("failed to scan sequence row: %w", err)
		}

		sequences = append(sequences, seq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sequence rows: %w", err)
	}

	return sequences, nil
}

func querySchemaObjects[T any](ctx context.Context, db *sql.DB, query string, schemas []string,
	newItem func() T, scanFields func(*T) []interface{}, objectType string) ([]T, error) {
	args := make([]interface{}, len(schemas))
	for i, schema := range schemas {
		args[i] = schema
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", objectType, err)
	}
	defer rows.Close()

	var items []T

	for rows.Next() {
		item := newItem()

		if err := rows.Scan(scanFields(&item)...); err != nil {
			return nil, fmt.Errorf("failed to scan %s row: %w", objectType, err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating %s rows: %w", objectType, err)
	}

	return items, nil
}

func (a *PostgreSQLAnalyzer) buildSchemaFilterQuery(baseQuery, schemaColumn, orderBy string, schemas []string) string {
	whereClause := " AND "

	switch len(schemas) {
	case 0:
		whereClause += fmt.Sprintf("NOT %s IN ('pg_catalog', 'information_schema', 'pg_toast')", schemaColumn)
	case 1:
		whereClause += fmt.Sprintf("%s = $1", schemaColumn)
	default:
		placeholders := make([]string, len(schemas))
		for i := range schemas {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

		whereClause += fmt.Sprintf("%s IN (%s)", schemaColumn, strings.Join(placeholders, ", "))
	}

	return baseQuery + whereClause + fmt.Sprintf(" ORDER BY %s", orderBy)
}

func (a *PostgreSQLAnalyzer) buildSequenceQuery(schemas []string) string {
	return a.buildSchemaFilterQuery(
		`SELECT schemaname AS schema_name, sequencename AS sequence_name, data_type, start_value, min_value, max_value, increment_by FROM pg_catalog.pg_sequences WHERE true`,
		"schemaname",
		"schemaname, sequencename",
		schemas,
	)
}
