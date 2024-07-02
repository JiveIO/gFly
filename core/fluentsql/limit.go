package fluentsql

import "fmt"

// Limit clause
type Limit struct {
	Limit  int
	Offset int
}

func (l *Limit) String() string {
	if l.Limit > 0 || l.Offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", l.Limit, l.Offset)
	}

	return ""
}
