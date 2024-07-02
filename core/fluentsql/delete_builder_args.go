package fluentsql

import (
	"fmt"
	"strings"
)

// Sql Get Query statement and Arguments
func (db *DeleteBuilder) Sql() (string, []any, error) {
	var args []any

	return db.StringArgs(args)
}

func (db *DeleteBuilder) StringArgs(args []any) (string, []any, error) {
	var queryParts []string
	var sqlStr string

	sqlStr, args = db.deleteStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	sqlStr, args = db.whereStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = db.orderByStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = db.limitStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sql := strings.Join(queryParts, " ")

	return sql, args, nil
}

func (u *Delete) StringArgs(args []any) (string, []any) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", u.Table))

	if u.Alias != "" {
		sb.WriteString(" " + u.Alias)
	}

	return sb.String(), args
}
