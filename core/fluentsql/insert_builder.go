package fluentsql

import (
	"strings"
)

// ===========================================================================================================
//										Insert Builder :: Structure
// ===========================================================================================================

// InsertBuilder struct

/*
INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    [(col_name [, col_name] ...)]
    { {VALUES | VALUE} (value_list) [, (value_list)] ... }
    [AS row_alias[(col_alias [, col_alias] ...)]]
    [ON DUPLICATE KEY UPDATE assignment_list]

INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    SET assignment_list
    [AS row_alias[(col_alias [, col_alias] ...)]]
    [ON DUPLICATE KEY UPDATE assignment_list]

INSERT [LOW_PRIORITY | HIGH_PRIORITY] [IGNORE]
    [INTO] tbl_name
    [PARTITION (partition_name [, partition_name] ...)]
    [(col_name [, col_name] ...)]
    { SELECT ...
      | TABLE table_name
      | VALUES row_constructor_list
    }
    [ON DUPLICATE KEY UPDATE assignment_list]

value:
    {expr | DEFAULT}

value_list:
    value [, value] ...

row_constructor_list:
    ROW(value_list)[, ROW(value_list)][, ...]

assignment:
    col_name =
          value
        | [row_alias.]col_name
        | [tbl_name.]col_name
        | [row_alias.]col_alias

assignment_list:
    assignment [, assignment] ...
*/

type InsertBuilder struct {
	insertStatement Insert
	rowStatement    InsertRows
	queryStatement  InsertQuery
}

// InsertInstance Insert Builder constructor
func InsertInstance() *InsertBuilder {
	return &InsertBuilder{}
}

// ===========================================================================================================
//										Insert Builder :: Operators
// ===========================================================================================================

func (ib *InsertBuilder) String() string {
	var queryParts []string
	var sqlStr string

	queryParts = append(queryParts, ib.insertStatement.String())

	sqlStr = ib.rowStatement.String()
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr = ib.queryStatement.String()
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sql := strings.Join(queryParts, " ")

	return sql
}

// Insert builder
func (ib *InsertBuilder) Insert(table string, columns ...string) *InsertBuilder {
	ib.insertStatement.Table = table
	ib.insertStatement.Columns = columns

	return ib
}

// Row builder
func (ib *InsertBuilder) Row(values ...any) *InsertBuilder {
	ib.rowStatement.Append(values...)

	return ib
}

// Query builder
func (ib *InsertBuilder) Query(query *QueryBuilder) *InsertBuilder {
	ib.queryStatement.Query = query

	return ib
}
