package fluentsql

import (
	"fmt"
	"strings"
)

// Delete clause
type Delete struct {
	Table any
	Alias string
}

func (u *Delete) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", u.Table))

	if u.Alias != "" {
		sb.WriteString(" " + u.Alias)
	}

	return sb.String()
}
