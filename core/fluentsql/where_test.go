package fluentsql

import (
	"testing"
)

// TestWhereBasic
func TestWhereBasic(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE salary > 1400": {
			Field: "salary",
			Opt:   Greater,
			Value: 1400,
		},
		"WHERE department_id = 5": {
			Field: "department_id",
			Opt:   Eq,
			Value: 5,
		},
		"WHERE department_id <> 8": {
			Field: "department_id",
			Opt:   NotEq,
			Value: 8,
		},
		"WHERE last_name = 'Chen'": {
			Field: "last_name",
			Opt:   Eq,
			Value: "Chen",
		},
		"WHERE hire_date >= '1999-01-01'": {
			Field: "hire_date",
			Opt:   GrEq,
			Value: "1999-01-01",
		},
		"WHERE DATE_PART('year', year) = 1999": { // PostgreSQL
			Field: FieldYear("year"),
			Opt:   Eq,
			Value: 1999,
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereNull
func TestWhereNull(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE salary IS NULL": {
			Field: "salary",
			Opt:   Null,
			Value: nil,
		},
		"WHERE phone_number IS NOT NULL": {
			Field: "phone_number",
			Opt:   NotNull,
			Value: nil,
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereLike
func TestWhereLike(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE first_name LIKE 'S%'": {
			Field: "first_name",
			Opt:   Like,
			Value: "S%",
		},
		"WHERE first_name NOT LIKE 'Sh%'": {
			Field: "first_name",
			Opt:   NotLike,
			Value: "Sh%",
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereIn
func TestWhereIn(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE job_id IN (8, 9, 10)": {
			Field: "job_id",
			Opt:   In,
			Value: []int{8, 9, 10},
		},
		"WHERE job_id NOT IN (7, 8, 9)": {
			Field: "job_id",
			Opt:   NotIn,
			Value: []int{7, 8, 9},
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereOr
func TestWhereOr(t *testing.T) {
	var conditions []Condition

	conditions = append(conditions, Condition{
		Field: "salary",
		Opt:   Eq,
		Value: 7000,
	}, Condition{
		Field: "salary",
		Opt:   Eq,
		Value: 8000,
		AndOr: Or,
	})

	testCases := map[string][]Condition{
		"WHERE salary = 7000 OR salary = 8000": conditions,
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition...)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereBetween
func TestWhereBetween(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE salary BETWEEN 9000 AND 12000": {
			Field: "salary",
			Opt:   Between,
			Value: ValueBetween{
				Low:  9000,
				High: 12000,
			},
		},
		"WHERE DATE_PART('year', hire_date) BETWEEN 1990 AND 1993": {
			Field: FieldYear("hire_date"),
			Opt:   Between,
			Value: ValueBetween{
				Low:  1990,
				High: 1993,
			},
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereSubquery
func TestWhereSubquery(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE salary >= ALL (SELECT salary FROM employees WHERE department_id = 8)": {
			Field: "salary",
			Opt:   GrEqAll,
			Value: QueryInstance().
				Select("salary").
				From("employees").
				Where("department_id", Eq, 8),
		},
		"WHERE salary > ANY (SELECT AVG(salary) FROM employees GROUP BY department_id)": {
			Field: "salary",
			Opt:   GreaterAny,
			Value: QueryInstance().
				Select("AVG(salary)").
				From("employees").
				GroupBy("department_id"),
		},
		"WHERE  EXISTS (SELECT 1 FROM dependents d WHERE d.employee_id = e.employee_id)": {
			Field: FieldEmpty(""),
			Opt:   Exists,
			Value: QueryInstance().
				Select("1").
				From("dependents", "d").
				Where("d.employee_id", Eq, ValueField("e.employee_id")),
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}

// TestWhereOthers
func TestWhereOthers(t *testing.T) {
	testCases := map[string]Condition{
		"WHERE NOT salary > 5000": {
			Field: FieldNot("salary"),
			Opt:   Greater,
			Value: 5000,
		},
		"WHERE department_id IN (SELECT department_id FROM departments WHERE department_name = 'Marketing' OR department_name = 'Sales')": {
			Field: "department_id",
			Opt:   In,
			Value: QueryInstance().
				Select("department_id").
				From("departments").
				Where("department_name", Eq, "Marketing").
				WhereOr("department_name", Eq, "Sales"),
		},
		"WHERE  NOT EXISTS (SELECT employee_id FROM dependents d WHERE d.employee_id = e.employee_id)": {
			Field: FieldEmpty(""),
			Opt:   NotExists,
			Value: QueryInstance().
				Select("employee_id").
				From("dependents", "d").
				Where("d.employee_id", Eq, ValueField("e.employee_id")),
		},
		"WHERE department_id IN (SELECT department_id FROM departments WHERE active = true AND (department_name = 'Technical' OR department_name = 'Sales'))": {
			Field: "department_id",
			Opt:   In,
			Value: QueryInstance().
				Select("department_id").
				From("departments").
				Where("active", Eq, true).
				WhereGroup(func(whereBuilder WhereBuilder) *WhereBuilder {
					whereBuilder.Where("department_name", Eq, "Technical").
						WhereOr("department_name", Eq, "Sales")

					return &whereBuilder
				}),
		},
	}

	for expected, condition := range testCases {
		whereTest := new(Where)
		whereTest.Append(condition)

		if whereTest.String() != expected {
			t.Fatalf(`Query %s != %s`, whereTest.String(), expected)
		}
	}
}
