package fluentsql

import (
	"fmt"
	"strings"
)

// GroupBy clause
type GroupBy struct {
	Items []string
}

func (g *GroupBy) Append(field ...string) {
	g.Items = append(g.Items, field...)
}

func (g *GroupBy) String() string {
	if len(g.Items) == 0 {
		return ""
	}

	return fmt.Sprintf("GROUP BY %s", strings.Join(g.Items, ", "))
}
