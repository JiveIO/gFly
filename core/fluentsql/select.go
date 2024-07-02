package fluentsql

import (
	"fmt"
	"strings"
)

// Select clause
type Select struct {
	// Columns type string or a QueryBuilder
	Columns []any
}

func (s *Select) String() string {
	selectOf := "*"

	if len(s.Columns) > 0 {
		var columns []string

		for _, col := range s.Columns {
			if valueCase, ok := col.(*Case); ok { // Column type string
				columns = append(columns, valueCase.String())
			} else if valueString, ok := col.(string); ok { // Column type string
				columns = append(columns, valueString)
			} else if valueFieldYear, ok := col.(FieldYear); ok { // Column type FieldYear
				columns = append(columns, valueFieldYear.String())
			} else if valueQueryBuilder, ok := col.(*QueryBuilder); ok { // Column type is QueryBuilder.
				selectQuery := valueQueryBuilder.String()

				if valueQueryBuilder.alias == "" {
					selectQuery = fmt.Sprintf("(%s)", valueQueryBuilder)
				}

				columns = append(columns, selectQuery)
			}
		}

		selectOf = strings.Join(columns, ", ")
	}

	return fmt.Sprintf("SELECT %s", selectOf)
}
