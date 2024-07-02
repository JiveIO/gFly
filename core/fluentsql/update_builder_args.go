package fluentsql

import (
	"fmt"
	"strings"
)

// Sql Get Query statement and Arguments
func (ub *UpdateBuilder) Sql() (string, []any, interface{}) {
	return ub.StringArgs()
}

func (ub *UpdateBuilder) StringArgs() (string, []any, error) {
	var queryParts []string
	var sql string
	var args []any

	sql, args = ub.updateStatement.StringArgs(args)
	queryParts = append(queryParts, sql)

	sql, args = ub.setStatement.StringArgs(args)
	queryParts = append(queryParts, sql)

	sql, args = ub.whereStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	sql, args = ub.orderByStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	sql, args = ub.limitStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	sql = strings.Join(queryParts, " ")

	return sql, args, nil
}

func (u *Update) StringArgs(args []any) (string, []any) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE %s", u.Table))

	if u.Alias != "" {
		sb.WriteString(" " + u.Alias)
	}

	return sb.String(), args
}

func (s *UpdateSet) StringArgs(args []any) (string, []any) {
	var setColumns []string

	for _, item := range s.Items {
		var sql string

		sql, args = item.StringArgs(args)

		setColumns = append(setColumns, sql)
	}

	return fmt.Sprintf("SET %s", strings.Join(setColumns, ", ")), args
}

func (s *UpdateItem) StringArgs(args []any) (string, []any) {
	// SET (field1, field2,...) = (int, string, ValueField...)
	// SET (field1, field2,...) = (SELECT * FROM table_name)
	if fieldStringSlice, ok := s.Field.([]string); ok {
		fieldStr := joinSlice(fieldStringSlice, ",")

		if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Value type is QueryBuilder.
			var _sql string
			_sql, args, _ = valueQueryBuilder.StringArgs(args)

			return fmt.Sprintf("(%s) = (%s)", fieldStr, _sql), args
		}

		if fieldAnySlice, ok := s.Value.([]any); ok { // Value type is slice
			var values []string
			for _, fieldAny := range fieldAnySlice {
				if valueField, ok := fieldAny.(ValueField); ok { // Value type is string.
					values = append(values, valueField.String())
				} else if valueString, ok := fieldAny.(string); ok { // Value type is string.
					args = append(args, valueString)
					valueStr := p(args)

					values = append(values, valueStr)
				} else { // Value type is int or float.
					args = append(args, fieldAny)
					valueStr := p(args)

					values = append(values, valueStr)
				}
			}

			valueStr := strings.Join(values, ", ")

			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueStr), args
		}

		return "", args
	}

	if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Value type is QueryBuilder.
		var _sql string
		_sql, args, _ = valueQueryBuilder.StringArgs(args)

		return fmt.Sprintf("%s = (%s)", s.Field, _sql), args
	}

	if valueField, ok := s.Value.(ValueField); ok { // Value type is string.
		return fmt.Sprintf("%s = %s", s.Field, valueField), args
	}

	if valueString, ok := s.Value.(string); ok { // Value type is string.
		args = append(args, valueString)
		valueStr := p(args)

		return fmt.Sprintf("%s = %s", s.Field, valueStr), args
	}
	args = append(args, s.Value)
	valueStr := p(args)

	return fmt.Sprintf("%s = %s", s.Field, valueStr), args
}
