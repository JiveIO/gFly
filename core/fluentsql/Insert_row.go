package fluentsql

import (
	"fmt"
	"strings"
)

type InsertRows struct {
	Rows []InsertRow
}

func (r *InsertRows) Append(values ...any) {
	r.Rows = append(r.Rows, InsertRow{Values: values})
}

func (r *InsertRows) String() string {
	var rowsStr []string

	for _, row := range r.Rows {
		rowsStr = append(rowsStr, row.String())
	}

	if len(rowsStr) == 0 {
		return ""
	}

	return fmt.Sprintf("VALUES %s", strings.Join(rowsStr, ", "))
}

type InsertRow struct {
	Values []any
}

func (ir *InsertRow) String() string {
	var rowStr []string

	for _, col := range ir.Values {
		if colField, ok := col.(ValueField); ok { // Value type is ValueField.
			rowStr = append(rowStr, colField.String())
		} else if colString, ok := col.(string); ok { // Value type is string.
			rowStr = append(rowStr, "'"+colString+"'")
		} else { // Value type is int or float.
			rowStr = append(rowStr, fmt.Sprintf("%v", col))
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(rowStr, ", "))
}
