package fluentsql

import (
	"fmt"
	"strings"
)

type OrderByDir int

const (
	Asc OrderByDir = iota
	Desc
)

type SortItem struct {
	Field     string
	Direction OrderByDir
}

// OrderBy clause
type OrderBy struct {
	Items []SortItem
}

func (o *SortItem) Dir() string {
	var sign string

	switch o.Direction {
	case Asc:
		sign = "ASC"
	case Desc:
		sign = "DESC"
	}

	return sign
}

func (o *OrderBy) Append(field string, dir OrderByDir) {
	o.Items = append(o.Items, SortItem{
		Field:     field,
		Direction: dir,
	})
}

func (o *OrderBy) String() string {
	if len(o.Items) == 0 {
		return ""
	}

	var orderItems []string
	for _, item := range o.Items {
		orderItems = append(orderItems, fmt.Sprintf("%s %s", item.Field, item.Dir()))
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderItems, ", "))
}
