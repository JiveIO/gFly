package fluentsql

import "strings"

// ===========================================================================================================
//										Update Builder :: Structure
// ===========================================================================================================

// UpdateBuilder struct
/*
UPDATE [LOW_PRIORITY] [IGNORE] table_reference
    SET assignment_list
    [WHERE where_condition]
    [ORDER BY ...]
    [LIMIT row_count]

value:
    {expr | DEFAULT}

assignment:
    col_name = value

assignment_list:
    assignment [, assignment] ...
*/
type UpdateBuilder struct {
	updateStatement  Update
	setStatement     UpdateSet
	whereStatement   Where
	orderByStatement OrderBy
	limitStatement   Limit
}

// UpdateInstance Update Builder constructor
func UpdateInstance() *UpdateBuilder {
	return &UpdateBuilder{}
}

// ===========================================================================================================
//										Update Builder :: Operators
// ===========================================================================================================

func (ub *UpdateBuilder) String() string {
	var queryParts []string

	queryParts = append(queryParts, ub.updateStatement.String(), ub.setStatement.String())

	whereSql := ub.whereStatement.String()
	if whereSql != "" {
		queryParts = append(queryParts, whereSql)
	}

	orderBySql := ub.orderByStatement.String()
	if orderBySql != "" {
		queryParts = append(queryParts, orderBySql)
	}

	limitSql := ub.limitStatement.String()
	if limitSql != "" {
		queryParts = append(queryParts, limitSql)
	}

	sql := strings.Join(queryParts, " ")

	return sql
}

// Update builder
func (ub *UpdateBuilder) Update(table any, alias ...string) *UpdateBuilder {
	ub.updateStatement.Table = table

	// Table alias
	if len(alias) > 0 {
		ub.updateStatement.Alias = alias[0]
	}

	return ub
}

func (ub *UpdateBuilder) Set(field, value any) *UpdateBuilder {
	ub.setStatement.Append(field, value)

	return ub
}

// Where builder
func (ub *UpdateBuilder) Where(field any, opt WhereOpt, value any) *UpdateBuilder {
	ub.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
	})

	return ub
}

// WhereOr builder
func (ub *UpdateBuilder) WhereOr(field any, opt WhereOpt, value any) *UpdateBuilder {
	ub.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: Or,
	})

	return ub
}

// WhereGroup combine multi where conditions into a group.
func (ub *UpdateBuilder) WhereGroup(groupCondition FnWhereBuilder) *UpdateBuilder {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*WhereInstance())

	cond := Condition{
		Group: whereBuilder.whereStatement.Conditions,
	}

	ub.whereStatement.Conditions = append(ub.whereStatement.Conditions, cond)

	return ub
}
