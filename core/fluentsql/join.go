package fluentsql

import (
	"fmt"
	"strings"
)

type JoinType int

const (
	InnerJoin JoinType = iota
	LeftJoin
	RightJoin
	FullOuterJoin
	CrossJoin
)

// JoinItem type struct
type JoinItem struct {
	Join      JoinType
	Table     string
	Condition Condition
}

func (j *JoinItem) opt() string {
	var sign string

	switch j.Join {
	case InnerJoin:
		sign = "INNER JOIN"
	case LeftJoin:
		sign = "LEFT JOIN"
	case RightJoin:
		sign = "RIGHT JOIN"
	case FullOuterJoin:
		sign = "FULL OUTER JOIN"
	case CrossJoin:
		sign = "CROSS JOIN"
	}

	return sign
}

// Join clause
type Join struct {
	Items []JoinItem
}

func (j *Join) Append(item JoinItem) {
	j.Items = append(j.Items, item)
}

func (j *Join) String() string {
	if len(j.Items) == 0 {
		return ""
	}

	var joinItems []string
	for _, item := range j.Items {
		joinStr := fmt.Sprintf("%s %s ON %s", item.opt(), item.Table, item.Condition.String())

		if item.Join == CrossJoin {
			joinStr = fmt.Sprintf("%s %s", item.opt(), item.Table)
		}

		joinItems = append(joinItems, joinStr)
	}

	return strings.Join(joinItems, " ")
}
