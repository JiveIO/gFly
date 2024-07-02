package fluentsql

import (
	"fmt"
	"reflect"
	"strings"
)

// Where clause
type Where struct {
	Conditions []Condition
}

func (w *Where) Append(conditions ...Condition) {
	w.Conditions = append(w.Conditions, conditions...)
}

func (w *Where) String() string {
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

	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
}

// Condition type struct
type Condition struct {
	// Field name of column type string | FieldNot
	Field any
	// Opt condition operators =, <>, >, <, >=, <=, LIKE, IN, NOT IN, BETWEEN
	Opt WhereOpt
	// Value data of condition
	Value any
	// AndOr Combination type with previous condition AND, OR. Default is AND
	AndOr WhereAndOr
	// Group sub-conditions in parentheses `()`.
	Group []Condition
}

type WhereOpt int

const (
	Eq         WhereOpt = iota // Comparison Operators
	NotEq                      // Comparison Operators
	Diff                       // Comparison Operators
	Greater                    // Comparison Operators
	Lesser                     // Comparison Operators
	GrEq                       // Comparison Operators
	LeEq                       // Comparison Operators
	Like                       // Logical Operators
	NotLike                    // Logical Operators
	In                         // Logical Operators
	NotIn                      // Logical Operators
	Between                    // Logical Operators
	NotBetween                 // Logical Operators
	Null                       // Logical Operators
	NotNull                    // Logical Operators
	Exists                     // Logical Operators
	NotExists                  // Logical Operators
	EqAny                      // Comparison + Logical Operators
	NotEqAny                   // Comparison + Logical Operators
	DiffAny                    // Comparison + Logical Operators
	GreaterAny                 // Comparison + Logical Operators
	LesserAny                  // Comparison + Logical Operators
	GrEqAny                    // Comparison + Logical Operators
	LeEqAny                    // Comparison + Logical Operators
	EqAll                      // Comparison + Logical Operators
	NotEqAll                   // Comparison + Logical Operators
	DiffAll                    // Comparison + Logical Operators
	GreaterAll                 // Comparison + Logical Operators
	LesserAll                  // Comparison + Logical Operators
	GrEqAll                    // Comparison + Logical Operators
	LeEqAll                    // Comparison + Logical Operators
)

func (c *Condition) opt() string {
	var sign string

	switch c.Opt {
	case Eq:
		sign = "="
	case NotEq:
		sign = "<>"
	case Diff:
		sign = "!="
	case Greater:
		sign = ">"
	case Lesser:
		sign = "<"
	case GrEq:
		sign = ">="
	case LeEq:
		sign = "<="
	case Like:
		sign = "LIKE"
	case NotLike:
		sign = "NOT LIKE"
	case In:
		sign = "IN"
	case NotIn:
		sign = "NOT IN"
	case Between:
		sign = "BETWEEN"
	case NotBetween:
		sign = "NOT BETWEEN"
	case Null:
		sign = "IS NULL"
	case NotNull:
		sign = "IS NOT NULL"
	case Exists:
		sign = "EXISTS"
	case NotExists:
		sign = "NOT EXISTS"
	case EqAny:
		sign = "= ANY"
	case NotEqAny:
		sign = "<> ANY"
	case DiffAny:
		sign = "!= ANY"
	case GreaterAny:
		sign = "> ANY"
	case LesserAny:
		sign = "< ANY"
	case GrEqAny:
		sign = ">= ANY"
	case LeEqAny:
		sign = "<= ANY"
	case EqAll:
		sign = "= ALL"
	case NotEqAll:
		sign = "<> ALL"
	case DiffAll:
		sign = "!= ALL"
	case GreaterAll:
		sign = "> ALL"
	case LesserAll:
		sign = "< ALL"
	case GrEqAll:
		sign = ">= ALL"
	case LeEqAll:
		sign = "<= ALL"
	}

	return sign
}

func (c *Condition) String() string {
	// Handle group conditions WhereGroup(groupCondition FnWhereBuilder)
	if len(c.Group) > 0 {
		var conditions []string

		for _, cond := range c.Group {
			var _condition = cond.String()

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
			return ""
		}

		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND "))
	}

	// WHERE Address IS NULL
	// WHERE Address IS NOT NULL
	if c.Opt == Null || c.Opt == NotNull {
		return fmt.Sprintf("%s %s", c.Field, c.opt())
	}

	// WHERE Country IN ('Germany', 'France', 'UK')
	// WHERE Age NOT IN (12, 31, 21)
	if c.Opt == In || c.Opt == NotIn {
		// Type of value
		typ := reflect.TypeOf(c.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			valuesStr := ""
			if values, ok := c.Value.([]string); ok {
				valuesStr = "'" + strings.Join(values, "', '") + "'"
			}
			if values, ok := c.Value.([]int); ok {
				valuesStr = joinSlice(values, ",")
			}

			return fmt.Sprintf("%s %s (%s)", c.Field, c.opt(), valuesStr)
		}
	}

	// WHERE Price BETWEEN 10 AND 20;
	// WHERE ProductName BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price NOT BETWEEN 10 AND 20;
	// WHERE ProductName NOT BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price BETWEEN 10 AND 20
	if c.Opt == Between || c.Opt == NotBetween {
		return fmt.Sprintf("%s %s %v", c.Field, c.opt(), c.Value)
	}

	// WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT 1 , 1);
	// WHERE CustomerID IN (SELECT CustomerID FROM Orders);
	// WHERE CustomerID NOT IN (SELECT CustomerID FROM Orders);
	// WHERE EXISTS (SELECT ProductName FROM Products);
	// WHERE NOT EXISTS (SELECT ProductName FROM Products);
	// WHERE ProductID = ANY (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	// WHERE ProductID > ALL (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	if valueQueryBuilder, ok := c.Value.(*QueryBuilder); ok { // Column type is a complex query.
		return fmt.Sprintf("%s %s (%v)", c.Field, c.opt(), valueQueryBuilder)
	}

	if valueString, ok := c.Value.(string); ok {
		return fmt.Sprintf("%s %s '%v'", c.Field, c.opt(), valueString)
	}

	return fmt.Sprintf("%s %s %v", c.Field, c.opt(), c.Value)
}

type WhereAndOr int

const (
	And WhereAndOr = iota
	Or
)

// ValueBetween for WhereOpt.Between or WhereOpt.NotBetween
type ValueBetween struct {
	Low  any
	High any
}

func (v ValueBetween) String() string {
	if _, ok := v.Low.(string); ok {
		// hire_date BETWEEN '1999-01-01' AND '2000-12-31'
		return fmt.Sprintf("'%v' AND '%v'", v.Low, v.High)
	}

	//  salary NOT BETWEEN 2500 AND 2900
	return fmt.Sprintf("%v AND %v", v.Low, v.High)
}

// ValueField to handle condition `WHERE c.column <WhereOpt> c.column_1`
//
// So, When build condition Where("d.employee_id", qb.Eq, qb.ValueField("e.employee_id")) to keep SQL string as
// d.employee_id = e.employee_id instead of
// d.employee_id = 'e.employee_id'
type ValueField string

func (v ValueField) String() string {
	return string(v)
}

// FieldNot to handle condition `WHERE NOT salary > 5000`
//
// So, When build condition Where(qb.FieldNot("salary"), qb.Greater, 5000) to keep SQL string as `NOT salary > 5000`
type FieldNot string

func (v FieldNot) String() string {
	return fmt.Sprintf("NOT %s", string(v))
}

// FieldEmpty to handle condition `WHERE NOT EXISTS (SELECT employee_id FROM dependents)`
type FieldEmpty string

func (v FieldEmpty) String() string {
	return string(v)
}

// FieldYear to handle conditions
// Where ("YEAR(hire_date) Between 1990 AND 1993", // MySQL
// Where ("DATE_PART('year', hire_date) Between 1990 AND 1993", // PostgreSQL
type FieldYear string

func (v FieldYear) String() string {
	if dbType == MySQL {
		return fmt.Sprintf("YEAR(%s)", string(v))
	} else if dbType == PostgreSQL {
		return fmt.Sprintf("DATE_PART('year', %s)", string(v))
	}

	// SQLite https://database.guide/how-to-extract-the-day-month-and-year-from-a-date-in-sqlite/
	return "strftime('%Y', " + string(v) + ")"
}
