package fluentsql

import (
	"fmt"
	"strings"
)

// Update clause
type Update struct {
	Table any
	Alias string
}

func (u *Update) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE %s", u.Table))

	if u.Alias != "" {
		sb.WriteString(" " + u.Alias)
	}

	return sb.String()
}

type UpdateItem struct {
	// Field name of column type string, []string
	Field any
	// Value data of string, int, ValueField, QueryBuilder, []any
	Value any
}

func (s *UpdateItem) String() string {
	// SET (field1, field2,...) = (int, string, ValueField...)
	// SET (field1, field2,...) = (SELECT * FROM table_name)
	if fieldStringSlice, ok := s.Field.([]string); ok {
		fieldStr := joinSlice(fieldStringSlice, ",")

		if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Value type is QueryBuilder.
			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueQueryBuilder)
		}

		if fieldAnySlice, ok := s.Value.([]any); ok { // Value type is slice
			var values []string
			for _, fieldAny := range fieldAnySlice {
				if valueField, ok := fieldAny.(ValueField); ok { // Value type is string.
					values = append(values, valueField.String())
				} else if valueString, ok := fieldAny.(string); ok { // Value type is string.
					values = append(values, "'"+valueString+"'")
				} else { // Value type is int or float.
					values = append(values, fmt.Sprintf("%v", fieldAny))
				}
			}

			valueStr := strings.Join(values, ", ")

			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueStr)
		}

		return ""
	}

	if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Value type is QueryBuilder.
		return fmt.Sprintf("%s = (%v)", s.Field, valueQueryBuilder)
	}

	if valueField, ok := s.Value.(ValueField); ok { // Value type is string.
		return fmt.Sprintf("%s = %s", s.Field, valueField)
	}

	if valueString, ok := s.Value.(string); ok { // Value type is string.
		return fmt.Sprintf("%s = '%s'", s.Field, valueString)
	}

	return fmt.Sprintf("%s = %v", s.Field, s.Value)
}

type UpdateSet struct {
	Items []UpdateItem
}

func (s *UpdateSet) Append(field, value any) {
	s.AppendItems(UpdateItem{
		Field: field,
		Value: value,
	})
}

func (s *UpdateSet) AppendItems(items ...UpdateItem) {
	s.Items = append(s.Items, items...)
}

func (s *UpdateSet) String() string {
	var setColumns []string

	for _, item := range s.Items {
		setColumns = append(setColumns, item.String())
	}

	return fmt.Sprintf("SET %s", strings.Join(setColumns, ", "))
}
