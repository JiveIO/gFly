package fluentsql

import (
	"fmt"
	"strings"
)

// Sql Get Query statement and Arguments
func (ib *InsertBuilder) Sql() (string, []any, error) {
	var args []any

	return ib.StringArgs(args)
}

func (ib *InsertBuilder) StringArgs(args []any) (string, []any, error) {
	var queryParts []string
	var sqlStr string

	sqlStr, args = ib.insertStatement.StringArgs(args)

	queryParts = append(queryParts, sqlStr)

	sqlStr, args = ib.rowStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = ib.queryStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sql := strings.Join(queryParts, " ")

	return sql, args, nil
}

func (i *Insert) StringArgs(args []any) (string, []any) {
	columnsStr := strings.Join(i.Columns, ", ")

	return fmt.Sprintf("INSERT INTO %s (%s)", i.Table, columnsStr), args
}

func (r *InsertRows) StringArgs(args []any) (string, []any) {
	var rowsStr []string
	var sqlStr string

	for _, row := range r.Rows {
		sqlStr, args = row.StringArgs(args)
		rowsStr = append(rowsStr, sqlStr)
	}

	if len(rowsStr) == 0 {
		return "", args
	}

	return fmt.Sprintf("VALUES %s", strings.Join(rowsStr, ", ")), args
}

func (ir *InsertRow) StringArgs(args []any) (string, []any) {
	var rowStr []string

	for _, col := range ir.Values {
		if colField, ok := col.(ValueField); ok { // Value type is ValueField.
			rowStr = append(rowStr, colField.String())
		} else if colString, ok := col.(string); ok { // Value type is string.
			args = append(args, colString)
			colStr := p(args)

			rowStr = append(rowStr, colStr)
		} else { // Value type is int or float.
			args = append(args, col)
			colStr := p(args)

			rowStr = append(rowStr, colStr)
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(rowStr, ", ")), args
}

func (q *InsertQuery) StringArgs(args []any) (string, []any) {
	if queryBuilder, ok := q.Query.(*QueryBuilder); ok {
		var sqlStr string

		sqlStr, args, _ = queryBuilder.StringArgs(args)

		return sqlStr, args
	}

	return "", args
}
