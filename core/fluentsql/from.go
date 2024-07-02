package fluentsql

import (
	"fmt"
	"strings"
)

// From clause
type From struct {
	Table any
	Alias string
}

func (f *From) String() string {
	var sb strings.Builder

	if _, ok := f.Table.(string); ok { // Table type string
		sb.WriteString(fmt.Sprintf("FROM %s", f.Table))
	} else if _, ok := f.Table.(*QueryBuilder); ok { // Table type is QueryBuilder.
		selectQuery := f.Table.(*QueryBuilder).String()

		if f.Table.(*QueryBuilder).alias == "" {
			sb.WriteString(fmt.Sprintf("FROM (%s)", selectQuery))
		} else {
			sb.WriteString(fmt.Sprintf("FROM %s", selectQuery))
		}
	}

	if f.Alias != "" {
		sb.WriteString(" " + f.Alias)
	}

	return sb.String()
}
