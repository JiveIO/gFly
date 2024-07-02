package fluentsql

type InsertQuery struct {
	Query any
}

func (q *InsertQuery) String() string {
	if queryBuilder, ok := q.Query.(*QueryBuilder); ok {
		return queryBuilder.String()
	}

	return ""
}
