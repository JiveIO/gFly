package fluentsql

// WhereBuilder struct
type WhereBuilder struct {
	whereStatement Where
}

// WhereInstance Query builder constructor
func WhereInstance() *WhereBuilder {
	return &WhereBuilder{}
}

// Where builder
func (wb *WhereBuilder) Where(field any, opt WhereOpt, value any) *WhereBuilder {
	wb.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
	})

	return wb
}

// WhereOr builder
func (wb *WhereBuilder) WhereOr(field any, opt WhereOpt, value any) *WhereBuilder {
	wb.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: Or,
	})

	return wb
}

// WhereGroup combine multi where conditions into a group.
func (wb *WhereBuilder) WhereGroup(groupCondition FnWhereBuilder) *WhereBuilder {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*WhereInstance())

	cond := Condition{
		Group: whereBuilder.whereStatement.Conditions,
	}

	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, cond)

	return wb
}

// WhereCondition appends multi conditions
func (wb *WhereBuilder) WhereCondition(conditions ...Condition) *WhereBuilder {
	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, conditions...)

	return wb
}

type FnWhereBuilder func(whereBuilder WhereBuilder) *WhereBuilder

func (wb *WhereBuilder) String() string {
	return wb.whereStatement.String()
}

func (wb *WhereBuilder) StringArgs(args []any) (string, []any) {
	return wb.whereStatement.StringArgs(args)
}

// Conditions explode conditions
func (wb *WhereBuilder) Conditions() []Condition {
	return wb.whereStatement.Conditions
}
