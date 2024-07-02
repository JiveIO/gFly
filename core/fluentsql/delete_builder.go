package fluentsql

import "strings"

// ===========================================================================================================
//										Delete Builder :: Structure
// ===========================================================================================================

// DeleteBuilder struct
/*
DELETE [LOW_PRIORITY] [QUICK] [IGNORE] FROM tbl_name [[AS] tbl_alias]
    [PARTITION (partition_name [, partition_name] ...)]
    [WHERE where_condition]
    [ORDER BY ...]
    [LIMIT row_count]
*/
type DeleteBuilder struct {
	deleteStatement  Delete
	whereStatement   Where
	orderByStatement OrderBy
	limitStatement   Limit
}

// DeleteInstance Delete Builder constructor
func DeleteInstance() *DeleteBuilder {
	return &DeleteBuilder{}
}

// ===========================================================================================================
//										Update Builder :: Operators
// ===========================================================================================================

func (db *DeleteBuilder) String() string {
	var queryParts []string

	queryParts = append(queryParts, db.deleteStatement.String())

	whereSql := db.whereStatement.String()
	if whereSql != "" {
		queryParts = append(queryParts, whereSql)
	}

	orderBySql := db.orderByStatement.String()
	if orderBySql != "" {
		queryParts = append(queryParts, orderBySql)
	}

	limitSql := db.limitStatement.String()
	if limitSql != "" {
		queryParts = append(queryParts, limitSql)
	}

	sql := strings.Join(queryParts, " ")

	return sql
}

// Delete builder
func (db *DeleteBuilder) Delete(table string, alias ...string) *DeleteBuilder {
	db.deleteStatement.Table = table

	if len(alias) > 0 {
		db.deleteStatement.Alias = alias[0]
	}

	return db
}

// Where builder
func (db *DeleteBuilder) Where(field any, opt WhereOpt, value any) *DeleteBuilder {
	db.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
	})

	return db
}

// WhereOr builder
func (db *DeleteBuilder) WhereOr(field any, opt WhereOpt, value any) *DeleteBuilder {
	db.whereStatement.Append(Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: Or,
	})

	return db
}

// WhereGroup combine multi where conditions into a group.
func (db *DeleteBuilder) WhereGroup(groupCondition FnWhereBuilder) *DeleteBuilder {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*WhereInstance())

	cond := Condition{
		Group: whereBuilder.whereStatement.Conditions,
	}

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, cond)

	return db
}
