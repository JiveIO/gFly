package db

import (
	"app/core/fluentsql"
	"database/sql"
	"fmt"
	"strings"
)

// ===========================================================================================================
// 											Fluent SQL - Query
// ===========================================================================================================

type FnQueryBuilder func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder

func (db *DB) FluentGet(builderFunc FnQueryBuilder, dest interface{}, selectOf ...string) error {
	selectFields := "*"
	if len(selectOf) > 0 {
		selectFields = strings.Join(selectOf, ",")
	}

	// Create Fluent builder
	sqlBuilder := fluentsql.Select(selectFields)
	sqlBuilder = sqlBuilder.PlaceholderFormat(fluentsql.Dollar)

	// Get Fluent SQL Builder from user
	sqlBuilder = builderFunc(sqlBuilder)

	// Adapts FluentSQL compatible with DB
	sqlStr, args, _ := sqlBuilder.ToSql()

	// Query data
	return db.Get(dest, sqlStr, args...)
}

func (db *DB) FluentSelect(builderFunc FnQueryBuilder, dest interface{}, total *int, selectOf ...string) error {
	selectFields := "*"
	if len(selectOf) > 0 {
		selectFields = strings.Join(selectOf, ",")
	}

	// Create Fluent builder
	sqlBuilder := fluentsql.Select(selectFields)
	sqlBuilder = sqlBuilder.PlaceholderFormat(fluentsql.Dollar)

	// Get Fluent SQL Builder from user
	sqlBuilder = builderFunc(sqlBuilder)

	// Adapts FluentSQL compatible with DB
	sqlStr, args, _ := sqlBuilder.ToSql()

	// Query COUNT
	err := db.count(sqlBuilder, total)
	if err != nil {
		return err
	}

	// Query data
	err = db.Select(dest, sqlStr, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) count(sqlBuilder fluentsql.SelectBuilder, total *int) error {
	// Build SQL without pagination
	sqlNoLimit, _, _ := sqlBuilder.RemoveOffset().RemoveLimit().ToSql()

	// Create COUNT query
	sqlBuilderCount := fluentsql.
		Select("COUNT(*) AS total").
		From(fmt.Sprintf("(%s) AS _result_out_", sqlNoLimit))

	// Extra SQL
	sqlStrCount, argsCount, _ := sqlBuilderCount.ToSql()

	return db.Get(total, sqlStrCount, argsCount...)
}

// ===========================================================================================================
// 											Fluent SQL - Insert
// ===========================================================================================================

type FnInsertBuilder func(query fluentsql.InsertBuilder) fluentsql.InsertBuilder

func (db *DB) FluentInsert(builderFunc FnInsertBuilder, table string) (sql.Result, error) {
	// Create Fluent builder
	sqlBuilder := fluentsql.Insert(table)
	sqlBuilder = sqlBuilder.PlaceholderFormat(fluentsql.Dollar)

	// Get Fluent SQL Builder from user
	sqlBuilder = builderFunc(sqlBuilder)

	// Adapts FluentSQL compatible with DB
	sqlStr, args, _ := sqlBuilder.ToSql()

	// Query data
	return db.Exec(sqlStr, args...)
}

// ===========================================================================================================
// 											Fluent SQL - Update
// ===========================================================================================================

type FnUpdateBuilder func(query fluentsql.UpdateBuilder) fluentsql.UpdateBuilder

func (db *DB) FluentUpdate(builderFunc FnUpdateBuilder, table string) (sql.Result, error) {
	// Create Fluent builder
	sqlBuilder := fluentsql.Update(table)
	sqlBuilder = sqlBuilder.PlaceholderFormat(fluentsql.Dollar)

	// Get Fluent SQL Builder from user
	sqlBuilder = builderFunc(sqlBuilder)

	// Adapts FluentSQL compatible with DB
	sqlStr, args, _ := sqlBuilder.ToSql()

	// Query data
	return db.Exec(sqlStr, args...)
}

// ===========================================================================================================
// 											Fluent SQL - Delete
// ===========================================================================================================

type FnDeleteBuilder func(query fluentsql.DeleteBuilder) fluentsql.DeleteBuilder

func (db *DB) FluentDelete(builderFunc FnDeleteBuilder, table string) (sql.Result, error) {
	// Create Fluent builder
	sqlBuilder := fluentsql.Delete(table)
	sqlBuilder = sqlBuilder.PlaceholderFormat(fluentsql.Dollar)

	// Get Fluent SQL Builder from user
	sqlBuilder = builderFunc(sqlBuilder)

	// Adapts FluentSQL compatible with DB
	sqlStr, args, _ := sqlBuilder.ToSql()

	// Query data
	return db.Exec(sqlStr, args...)
}
