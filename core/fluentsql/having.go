package fluentsql

import (
	"fmt"
	"strings"
)

// Having clause
type Having struct {
	Where
}

func (w *Having) String() string {
	var conditions []string

	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			var _condition = cond.String()

			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				last := len(conditions) - 1

				conditions[last] += _orCondition
			} else {
				conditions = append(conditions, _condition)
			}
		}
	}

	// No WHERE condition
	if len(conditions) == 0 {
		return ""
	}

	return fmt.Sprintf("HAVING %s", strings.Join(conditions, " AND "))
}
