package fluentsql

import (
	"fmt"
	"reflect"
	"strings"
)

// Sql Get Query statement and Arguments
func (qb *QueryBuilder) Sql() (string, []any, error) {
	var args []any

	return qb.StringArgs(args)
}

// StringArgs builder
func (qb *QueryBuilder) StringArgs(args []any) (string, []any, error) {
	var queryParts []string
	var sqlStr string

	sqlStr, args = qb.selectStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	sqlStr, args = qb.fromStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	sqlStr, args = qb.joinStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.whereStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.groupByStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.havingStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.orderByStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.limitStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.fetchStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr = strings.Join(queryParts, " ")

	if qb.alias != "" {
		sqlStr = fmt.Sprintf("(%s) AS %s",
			sqlStr,
			qb.alias)
	}

	return sqlStr, args, nil
}

func (s *Select) StringArgs(args []any) (string, []any) {
	selectOf := "*"

	if len(s.Columns) > 0 {
		var columns []string

		for _, col := range s.Columns {
			var sqlPart string
			if _, ok := col.(*Case); ok { // Column type string
				sqlPart, args = col.(*Case).StringArgs(args)
				columns = append(columns, sqlPart)
			} else if valueString, ok := col.(string); ok { // Column type string
				columns = append(columns, valueString)
			} else if valueFieldYear, ok := col.(FieldYear); ok { // Column type FieldYear
				columns = append(columns, valueFieldYear.String())
			} else if valueQueryBuilder, ok := col.(*QueryBuilder); ok { // Column type is QueryBuilder.
				var selectQuery string
				selectQuery, args, _ = valueQueryBuilder.StringArgs(args)

				if col.(*QueryBuilder).alias == "" {
					selectQuery = fmt.Sprintf("(%s)", selectQuery)
				}

				columns = append(columns, selectQuery)
			}
		}

		selectOf = strings.Join(columns, ", ")
	}

	return fmt.Sprintf("SELECT %s", selectOf), args
}

func (f *From) StringArgs(args []any) (string, []any) {
	var sb strings.Builder

	if _, ok := f.Table.(string); ok { // Table type string
		sb.WriteString(fmt.Sprintf("FROM %s", f.Table))
	} else if _, ok := f.Table.(*QueryBuilder); ok { // Table type is QueryBuilder.
		var selectQuery string
		selectQuery, args, _ = f.Table.(*QueryBuilder).StringArgs(args)

		if f.Table.(*QueryBuilder).alias == "" {
			sb.WriteString(fmt.Sprintf("FROM (%s)", selectQuery))
		} else {
			sb.WriteString(fmt.Sprintf("FROM %s", selectQuery))
		}
	}

	if f.Alias != "" {
		sb.WriteString(" " + f.Alias)
	}

	return sb.String(), args
}

func (j *Join) StringArgs(args []any) (string, []any) {
	if len(j.Items) == 0 {
		return "", args
	}

	var joinItems []string
	for _, item := range j.Items {
		var cond string
		cond, args = item.Condition.StringArgs(args)

		joinStr := fmt.Sprintf("%s %s ON %s", item.opt(), item.Table, cond)

		if item.Join == CrossJoin {
			joinStr = fmt.Sprintf("%s %s", item.opt(), item.Table)
		}

		joinItems = append(joinItems, joinStr)
	}

	return strings.Join(joinItems, " "), args
}

func (w *Where) StringArgs(args []any) (string, []any) {
	var conditions []string

	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			var _condition string
			_condition, args = cond.StringArgs(args)

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
		return "", args
	}

	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND ")), args
}

func (c *Condition) StringArgs(args []any) (string, []any) {
	// Handle group conditions WhereGroup(groupCondition FnWhereBuilder)
	if len(c.Group) > 0 {
		var conditions []string

		for _, cond := range c.Group {
			var _condition string
			_condition, args = cond.StringArgs(args)

			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				last := len(conditions) - 1

				// OR with previous condition
				conditions[last] += _orCondition
			} else {
				conditions = append(conditions, _condition)
			}
		}

		// No WHERE condition
		if len(conditions) == 0 {
			return "", args
		}

		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND ")), args
	}

	// Not include value type ValueField to arguments.
	if valueField, ok := c.Value.(ValueField); ok {
		return fmt.Sprintf("%s %s %s", c.Field, c.opt(), valueField), args
	}

	// WHERE Address IS NULL
	// WHERE Address IS NOT NULL
	if c.Opt == Null || c.Opt == NotNull {
		return fmt.Sprintf("%s %s", c.Field, c.opt()), args
	}

	// WHERE Country IN ('Germany', 'France', 'UK')
	// WHERE Age NOT IN (12, 31, 21)
	if c.Opt == In || c.Opt == NotIn {
		// Type of value
		typ := reflect.TypeOf(c.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			var valuesStr []string

			if values, ok := c.Value.([]string); ok {
				for _, val := range values {
					args = append(args, val)
					valuesStr = append(valuesStr, p(args))
				}
			}
			if values, ok := c.Value.([]int); ok {
				for _, val := range values {
					args = append(args, val)
					valuesStr = append(valuesStr, p(args))
				}
			}

			return fmt.Sprintf("%s %s (%s)", c.Field, c.opt(), strings.Join(valuesStr, ", ")), args
		}
	}

	// WHERE Price BETWEEN 10 AND 20;
	// WHERE ProductName BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price NOT BETWEEN 10 AND 20;
	// WHERE ProductName NOT BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price BETWEEN 10 AND 20
	if c.Opt == Between || c.Opt == NotBetween {
		var betweenValue string
		betweenValue, args = c.Value.(ValueBetween).StringArgs(args)

		return fmt.Sprintf("%s %s %v", c.Field, c.opt(), betweenValue), args
	}

	// WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT 1 , 1);
	// WHERE CustomerID IN (SELECT CustomerID FROM Orders);
	// WHERE CustomerID NOT IN (SELECT CustomerID FROM Orders);
	// WHERE EXISTS (SELECT ProductName FROM Products);
	// WHERE NOT EXISTS (SELECT ProductName FROM Products);
	// WHERE ProductID = ANY (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	// WHERE ProductID > ALL (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	if valueQueryBuilder, ok := c.Value.(*QueryBuilder); ok { // Column type is a complex query.
		var queryBuilderStr string
		queryBuilderStr, args, _ = valueQueryBuilder.StringArgs(args)

		return fmt.Sprintf("%s %s (%v)", c.Field, c.opt(), queryBuilderStr), args
	}

	if valueString, ok := c.Value.(string); ok {
		args = append(args, valueString)

		return fmt.Sprintf("%s %s %s", c.Field, c.opt(), p(args)), args
	}

	args = append(args, c.Value)

	return fmt.Sprintf("%s %s %s", c.Field, c.opt(), p(args)), args
}

func (v ValueBetween) StringArgs(args []any) (string, []any) {
	args = append(args, v.Low)
	pLow := p(args)
	args = append(args, v.High)
	pHigh := p(args)

	// hire_date BETWEEN '1999-01-01' AND '2000-12-31'
	// salary NOT BETWEEN 2500 AND 2900
	return fmt.Sprintf("%v AND %v", pLow, pHigh), args
}

func (v FieldYear) StringArgs(args []any) (string, []any) {
	args = append(args, string(v))
	value := p(args)

	if dbType == MySQL {
		return fmt.Sprintf("YEAR(%s)", value), args
	} else if dbType == PostgreSQL {
		return fmt.Sprintf("DATE_PART('year', %s)", value), args
	}

	// SQLite
	return "strftime('%Y', ?)", args
}

func (g *GroupBy) StringArgs(args []any) (string, []any) {
	if len(g.Items) == 0 {
		return "", args
	}

	return fmt.Sprintf("GROUP BY %s", strings.Join(g.Items, ", ")), args
}

func (w *Having) StringArgs(args []any) (string, []any) {
	var conditions []string

	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			var _condition string
			_condition, args = cond.StringArgs(args)

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
		return "", args
	}

	return fmt.Sprintf("HAVING %s", strings.Join(conditions, " AND ")), args
}

func (o *OrderBy) StringArgs(args []any) (string, []any) {
	if len(o.Items) == 0 {
		return "", args
	}

	var orderItems []string
	for _, item := range o.Items {
		orderItems = append(orderItems, fmt.Sprintf("%s %s", item.Field, item.Dir()))
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderItems, ", ")), args
}

func (l *Limit) StringArgs(args []any) (string, []any) {
	if l.Limit > 0 || l.Offset > 0 {
		args = append(args, l.Limit)
		pLimit := p(args)
		args = append(args, l.Offset)
		pOffset := p(args)

		return fmt.Sprintf("LIMIT %s OFFSET %s", pLimit, pOffset), args
	}

	return "", args
}

func (f *Fetch) StringArgs(args []any) (string, []any) {
	if f.Fetch > 0 || f.Offset > 0 {
		args = append(args, f.Offset)
		pOffset := p(args)
		args = append(args, f.Fetch)
		pFetch := p(args)

		return fmt.Sprintf("OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", pOffset, pFetch), args
	}

	return "", args
}

func (c *WhenCase) StringArgs(args []any) (string, []any) {
	args = append(args, c.Value)

	if valueConditions, ok := c.Conditions.([]Condition); ok {
		var cons []string
		var sqlPart string
		for _, condition := range valueConditions {
			sqlPart, args = condition.StringArgs(args)

			cons = append(cons, sqlPart)
		}

		return fmt.Sprintf("WHEN %s THEN '?'", strings.Join(cons, " AND ")), args
	}

	return fmt.Sprintf("WHEN %v THEN '?'", c.Conditions), args
}

func (c *Case) StringArgs(args []any) (string, []any) {
	var whenCases []string
	var sqlPart string

	for _, whenClause := range c.WhenClauses {
		sqlPart, args = whenClause.StringArgs(args)

		whenCases = append(whenCases, sqlPart)
	}

	return fmt.Sprintf("CASE %s %s END %s", c.Exp, strings.Join(whenCases, " "), c.Name), args
}
