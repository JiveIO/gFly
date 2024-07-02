package fluentsql

import (
	"testing"
)

// TestQueryBasic
func TestQueryBasic(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT first_name, last_name, salary, salary * 1.05 AS new_salary FROM employees": QueryInstance().
			Select("first_name", "last_name", "salary", "salary * 1.05 AS new_salary").
			From("employees"),
		"SELECT first_name, last_name, hire_date, salary FROM employees ORDER BY first_name ASC, last_name DESC": QueryInstance().
			Select("first_name", "last_name", "hire_date", "salary").
			From("employees").
			OrderBy("first_name", Asc).
			OrderBy("last_name", Desc),
		"SELECT DISTINCT salary FROM employees ORDER BY salary DESC": QueryInstance().
			Select("DISTINCT salary").
			From("employees").
			OrderBy("salary", Desc),
		"SELECT employee_id, first_name, last_name FROM employees ORDER BY first_name ASC LIMIT 5 OFFSET 3": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees").
			OrderBy("first_name", Asc).
			Limit(5, 3),
		"SELECT employee_id, first_name, last_name, hire_date FROM employees WHERE DATE_PART('year', hire_date) = 1999 ORDER BY hire_date DESC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "hire_date").
			From("employees").
			Where(FieldYear("hire_date"), Eq, 1999).
			OrderBy("hire_date", Desc),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary = 7000 OR salary = 8000 ORDER BY salary ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("salary", Eq, 7000).
			WhereOr("salary", Eq, 8000).
			OrderBy("salary", Asc),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryLikeAndInAndNotNull
func TestQueryLikeAndInAndNotNull(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT first_name, last_name, salary, salary * 1.05 AS new_salary FROM employees": QueryInstance().
			Select("first_name", "last_name", "salary", "salary * 1.05 AS new_salary").
			From("employees"),
		"SELECT employee_id, first_name, last_name FROM employees WHERE first_name LIKE 'Jo%' ORDER BY first_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees").
			Where("first_name", Like, "Jo%").
			OrderBy("first_name", Asc),
		"SELECT employee_id, first_name, last_name FROM employees WHERE first_name LIKE 'S%' AND first_name NOT LIKE 'Sh%' ORDER BY first_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees").
			Where("first_name", Like, "S%").
			Where("first_name", NotLike, "Sh%").
			OrderBy("first_name", Asc),
		"SELECT employee_id, first_name, last_name, department_id FROM employees WHERE department_id IN (8, 9) ORDER BY department_id ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "department_id").
			From("employees").
			Where("department_id", In, []int{8, 9}).
			OrderBy("department_id", Asc),
		"SELECT employee_id, first_name, last_name, phone_number FROM employees WHERE phone_number IS NOT NULL": QueryInstance().
			Select("employee_id", "first_name", "last_name", "phone_number").
			From("employees").
			Where("phone_number", NotNull, nil),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryFieldYear
func TestQueryFieldYear(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT employee_id, first_name, last_name, hire_date FROM employees WHERE DATE_PART('year', hire_date) = 1999 ORDER BY hire_date DESC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "hire_date").
			From("employees").
			Where(FieldYear("hire_date"), Eq, 1999).
			OrderBy("hire_date", Desc),
		"SELECT employee_id, first_name, last_name, DATE_PART('year', hire_date) joined_year FROM employees WHERE DATE_PART('year', hire_date) BETWEEN 1990 AND 1993 ORDER BY hire_date ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", FieldYear("hire_date").String()+" joined_year").
			From("employees").
			Where(FieldYear("hire_date"),
				Between, ValueBetween{
					Low:  1990,
					High: 1993,
				}).
			OrderBy("hire_date", Asc),
		"SELECT employee_id, first_name, last_name, DATE_PART('year', hire_date) FROM employees WHERE DATE_PART('year', hire_date) BETWEEN 1990 AND 1993 ORDER BY hire_date ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", FieldYear("hire_date")).
			From("employees").
			Where(FieldYear("hire_date"),
				Between, ValueBetween{
					Low:  1990,
					High: 1993,
				}).
			OrderBy("hire_date", Asc),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryLikeAndInAndNotNull
func TestQueryBetweenCase(t *testing.T) {
	var conditionsLow []Condition
	conditionsLow = append(conditionsLow, Condition{
		Field: "salary",
		Opt:   Lesser,
		Value: 3000,
	})

	var conditionsAverage []Condition
	conditionsAverage = append(conditionsAverage, Condition{
		Field: "salary",
		Opt:   GrEq,
		Value: 3000,
	}, Condition{
		Field: "salary",
		Opt:   LeEq,
		Value: 5000,
	})

	var conditionsHigh []Condition
	conditionsHigh = append(conditionsHigh, Condition{
		Field: "salary",
		Opt:   Greater,
		Value: 5000,
	})

	fieldCase := FieldCase("", "evaluation")
	fieldCase.When(conditionsLow, "Low")
	fieldCase.When(conditionsAverage, "Average")
	fieldCase.When(conditionsHigh, "High")

	testCases := map[string]*QueryBuilder{
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary NOT BETWEEN 3000 AND 5000": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			Where("salary", NotBetween, ValueBetween{Low: 3000, High: 5000}).
			From("employees"),
		"SELECT employee_id, first_name, last_name, DATE_PART('year', hire_date) joined_year FROM employees WHERE DATE_PART('year', hire_date) BETWEEN 1990 AND 1993 ORDER BY hire_date ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", FieldYear("hire_date").String()+" joined_year").
			From("employees").
			Where(FieldYear("hire_date"),
				Between, ValueBetween{
					Low:  1990,
					High: 1993,
				}).
			OrderBy("hire_date", Asc),
		"SELECT first_name, last_name, CASE  WHEN salary < 3000 THEN 'Low' WHEN salary >= 3000 AND salary <= 5000 THEN 'Average' WHEN salary > 5000 THEN 'High' END evaluation FROM employees": QueryInstance().
			Select("first_name", "last_name", fieldCase).
			From("employees"),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQuerySubquery
func TestQuerySubquery(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT 1 OFFSET 1)": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("salary", Eq,
				QueryInstance().
					Select("DISTINCT salary").
					From("employees").
					OrderBy("salary", Desc).
					Limit(1, 1),
			),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary >= ALL (SELECT salary FROM employees WHERE department_id = 8) ORDER BY salary DESC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("salary", GrEqAll,
				QueryInstance().
					Select("salary").
					From("employees").
					Where("department_id", Eq, 8),
			).
			OrderBy("salary", Desc),
		"SELECT employee_id, first_name, last_name, salary FROM employees e WHERE  EXISTS (SELECT 1 FROM dependents d WHERE d.employee_id = e.employee_id) ORDER BY first_name ASC, last_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees", "e").
			Where(FieldEmpty(""), Exists,
				QueryInstance().
					Select("1").
					From("dependents", "d").
					Where("d.employee_id", Eq, ValueField("e.employee_id")),
			).
			OrderBy("first_name", Asc).
			OrderBy("last_name", Asc),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE department_id IN (SELECT department_id FROM departments WHERE department_name = 'Marketing' OR department_name = 'Sales')": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("department_id", In,
				QueryInstance().
					Select("department_id").
					From("departments").
					Where("department_name", Eq, "Marketing").
					WhereOr("department_name", Eq, "Sales"),
			),
		"SELECT employee_id, first_name, last_name FROM employees e WHERE  NOT EXISTS (SELECT employee_id FROM dependents d WHERE d.employee_id = e.employee_id)": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees", "e").
			Where(FieldEmpty(""), NotExists,
				QueryInstance().
					Select("employee_id").
					From("dependents", "d").
					Where("d.employee_id", Eq, ValueField("e.employee_id")),
			),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryAlias
func TestQueryAlias(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT inv_no AS invoice_no, amount, due_date AS 'Due date', cust_no 'Customer No' FROM invoices": QueryInstance().
			Select("inv_no AS invoice_no", "amount", "due_date AS 'Due date'", "cust_no 'Customer No'").
			From("invoices"),
		"SELECT first_name, last_name, salary * 1.1 AS new_salary FROM employees WHERE new_salary > 5000": QueryInstance().
			Select("first_name", "last_name", "salary * 1.1 AS new_salary").
			From("employees").
			Where("new_salary", Greater, 5000),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryJoin
func TestQueryJoin(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT first_name, last_name, employees.department_id, departments.department_id, department_name FROM employees INNER JOIN departments ON departments.department_id = employees.department_id WHERE employees.department_id IN (1, 2, 3)": QueryInstance().
			Select("first_name", "last_name", "employees.department_id", "departments.department_id", "department_name").
			From("employees").
			Join(InnerJoin, "departments", Condition{
				Field: "departments.department_id",
				Opt:   Eq,
				Value: ValueField("employees.department_id"),
			}).
			Where("employees.department_id", In, []int{1, 2, 3}),
		"SELECT first_name, last_name, job_title, department_name FROM employees e INNER JOIN departments d ON d.department_id = e.department_id INNER JOIN jobs j ON j.job_id = e.job_id WHERE e.department_id IN (1, 2, 3)": QueryInstance().
			Select("first_name", "last_name", "job_title", "department_name").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			Join(InnerJoin, "jobs j", Condition{
				Field: "j.job_id",
				Opt:   Eq,
				Value: ValueField("e.job_id"),
			}).
			Where("e.department_id", In, []int{1, 2, 3}),
		"SELECT c.country_name, c.country_id, l.country_id, l.street_address, l.city FROM countries c LEFT JOIN locations l ON l.country_id = c.country_id WHERE c.country_id IN ('US', 'UK', 'CN')": QueryInstance().
			Select("c.country_name", "c.country_id", "l.country_id", "l.street_address", "l.city").
			From("countries", "c").
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("c.country_id", In, []string{"US", "UK", "CN"}),
		"SELECT country_name FROM countries c LEFT JOIN locations l ON l.country_id = c.country_id WHERE l.location_id IS NULL ORDER BY country_name ASC": QueryInstance().
			Select("country_name").
			From("countries", "c").
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("l.location_id", Null, nil).
			OrderBy("country_name", Asc),
		"SELECT r.region_name, c.country_name, l.street_address, l.city FROM regions r LEFT JOIN countries c ON c.region_id = r.region_id LEFT JOIN locations l ON l.country_id = c.country_id WHERE c.country_id IN ('US', 'UK', 'CN')": QueryInstance().
			Select("r.region_name", "c.country_name", "l.street_address", "l.city").
			From("regions", "r").
			Join(LeftJoin, "countries c", Condition{
				Field: "c.region_id",
				Opt:   Eq,
				Value: ValueField("r.region_id"),
			}).
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("c.country_id", In, []string{"US", "UK", "CN"}),
		"SELECT e.first_name || ' ' || e.last_name AS employee, m.first_name || ' ' || m.last_name AS manager FROM employees e INNER JOIN employees m ON m.employee_id = e.manager_id ORDER BY manager ASC": QueryInstance().
			Select("e.first_name || ' ' || e.last_name AS employee", "m.first_name || ' ' || m.last_name AS manager").
			From("employees", "e").
			Join(InnerJoin, "employees m", Condition{
				Field: "m.employee_id",
				Opt:   Eq,
				Value: ValueField("e.manager_id"),
			}).
			OrderBy("manager", Asc),
		"SELECT basket_name, fruit_name FROM fruits FULL OUTER JOIN baskets ON baskets.basket_id = fruits.basket_id WHERE fruit_name IS NULL": QueryInstance().
			Select("basket_name", "fruit_name").
			From("fruits").
			Join(FullOuterJoin, "baskets", Condition{
				Field: "baskets.basket_id",
				Opt:   Eq,
				Value: ValueField("fruits.basket_id"),
			}).
			Where("fruit_name", Null, nil),
		"SELECT sales_org, channel FROM sales_organization CROSS JOIN sales_channel": QueryInstance().
			Select("sales_org", "channel").
			From("sales_organization").
			Join(CrossJoin, "sales_channel", Condition{}),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryGroupByHaving
func TestQueryGroupByHaving(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT department_name, COUNT(employee_id) headcount FROM employees e INNER JOIN departments d ON d.department_id = e.department_id GROUP BY department_name": QueryInstance().
			Select("department_name", "COUNT(employee_id) headcount").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			GroupBy("department_name"),
		"SELECT department_name, COUNT(employee_id) headcount FROM employees e INNER JOIN departments d ON d.department_id = e.department_id GROUP BY department_name HAVING headcount > 5 ORDER BY headcount DESC": QueryInstance().
			Select("department_name", "COUNT(employee_id) headcount").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			GroupBy("department_name").
			Having("headcount", Greater, 5).
			OrderBy("headcount", Desc),
		"SELECT e.department_id, department_name, ROUND(AVG(salary), 2) FROM employees e INNER JOIN departments d ON d.department_id = e.department_id GROUP BY e.department_id HAVING AVG(salary) BETWEEN 5000 AND 7000 ORDER BY AVG(salary) ASC": QueryInstance().
			Select("e.department_id", "department_name", "ROUND(AVG(salary), 2)").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			GroupBy("e.department_id").
			Having("AVG(salary)", Between, ValueBetween{
				Low:  5000,
				High: 7000,
			}).
			OrderBy("AVG(salary)", Asc),
		"SELECT BillingDate, COUNT(*) AS BillingQty, SUM(BillingTotal) AS BillingSum FROM Billings WHERE BillingDate BETWEEN '2002-05-01' AND '2002-05-31' GROUP BY BillingDate HAVING COUNT(*) > 1 AND SUM(BillingTotal) > 100 ORDER BY BillingDate DESC": QueryInstance().
			Select("BillingDate", "COUNT(*) AS BillingQty", "SUM(BillingTotal) AS BillingSum").
			From("Billings").
			Where("BillingDate", Between, ValueBetween{
				Low:  "2002-05-01",
				High: "2002-05-31",
			}).
			GroupBy("BillingDate").
			Having("COUNT(*)", Greater, 1).
			Having("SUM(BillingTotal)", Greater, 100).
			OrderBy("BillingDate", Desc),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestQueryArgs
func TestQueryArgs(t *testing.T) {
	testCases := map[string]*QueryBuilder{
		"SELECT first_name, last_name FROM employees WHERE department_id IN ($1, $2, $3)": QueryInstance().
			Select("first_name", "last_name").
			From("employees").
			Where("department_id", In, []int{1, 2, 3}),
		"SELECT BillingDate, COUNT(*) AS BillingQty, SUM(BillingTotal) AS BillingSum FROM Billings WHERE BillingDate BETWEEN $1 AND $2 GROUP BY BillingDate HAVING COUNT(*) > $3 AND SUM(BillingTotal) > $4 ORDER BY BillingDate DESC": QueryInstance().
			Select("BillingDate", "COUNT(*) AS BillingQty", "SUM(BillingTotal) AS BillingSum").
			From("Billings").
			Where("BillingDate", Between, ValueBetween{
				Low:  "2002-05-01",
				High: "2002-05-31",
			}).
			GroupBy("BillingDate").
			Having("COUNT(*)", Greater, 1).
			Having("SUM(BillingTotal)", Greater, 100).
			OrderBy("BillingDate", Desc),
		"SELECT r.region_name, c.country_name, l.street_address, l.city FROM regions r LEFT JOIN countries c ON c.region_id = r.region_id LEFT JOIN locations l ON l.country_id = c.country_id WHERE c.country_id IN ($1, $2, $3)": QueryInstance().
			Select("r.region_name", "c.country_name", "l.street_address", "l.city").
			From("regions", "r").
			Join(LeftJoin, "countries c", Condition{
				Field: "c.region_id",
				Opt:   Eq,
				Value: ValueField("r.region_id"),
			}).
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("c.country_id", In, []string{"US", "UK", "CN"}),
		"SELECT department_name, COUNT(employee_id) headcount FROM employees e INNER JOIN departments d ON d.department_id = e.department_id GROUP BY department_name HAVING headcount > $1 ORDER BY headcount DESC": QueryInstance().
			Select("department_name", "COUNT(employee_id) headcount").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			GroupBy("department_name").
			Having("headcount", Greater, 5).
			OrderBy("headcount", Desc),
		"SELECT e.department_id, department_name, ROUND(AVG(salary), 2) FROM employees e INNER JOIN departments d ON d.department_id = e.department_id GROUP BY e.department_id HAVING AVG(salary) BETWEEN $1 AND $2 ORDER BY AVG(salary) ASC": QueryInstance().
			Select("e.department_id", "department_name", "ROUND(AVG(salary), 2)").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			GroupBy("e.department_id").
			Having("AVG(salary)", Between, ValueBetween{
				Low:  5000,
				High: 7000,
			}).
			OrderBy("AVG(salary)", Asc),
		"SELECT employee_id, first_name, last_name FROM employees WHERE first_name LIKE $1 ORDER BY first_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees").
			Where("first_name", Like, "Jo%").
			OrderBy("first_name", Asc),
		"SELECT employee_id, first_name, last_name FROM employees WHERE first_name LIKE $1 AND first_name NOT LIKE $2 ORDER BY first_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name").
			From("employees").
			Where("first_name", Like, "S%").
			Where("first_name", NotLike, "Sh%").
			OrderBy("first_name", Asc),
		"SELECT employee_id, first_name, last_name, department_id FROM employees WHERE department_id IN ($1, $2) ORDER BY department_id ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "department_id").
			From("employees").
			Where("department_id", In, []int{8, 9}).
			OrderBy("department_id", Asc),
		"SELECT employee_id, first_name, last_name, hire_date FROM employees WHERE DATE_PART('year', hire_date) = $1 ORDER BY hire_date DESC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "hire_date").
			From("employees").
			Where(FieldYear("hire_date"), Eq, 1999).
			OrderBy("hire_date", Desc),
		"SELECT employee_id, first_name, last_name, DATE_PART('year', hire_date) FROM employees WHERE DATE_PART('year', hire_date) BETWEEN $1 AND $2 ORDER BY hire_date ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", FieldYear("hire_date")).
			From("employees").
			Where(FieldYear("hire_date"),
				Between, ValueBetween{
					Low:  1990,
					High: 1993,
				}).
			OrderBy("hire_date", Asc),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary NOT BETWEEN $1 AND $2": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			Where("salary", NotBetween, ValueBetween{Low: 3000, High: 5000}).
			From("employees"),
		"SELECT employee_id, first_name, last_name, DATE_PART('year', hire_date) joined_year FROM employees WHERE DATE_PART('year', hire_date) BETWEEN $1 AND $2 ORDER BY hire_date ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", FieldYear("hire_date").String()+" joined_year").
			From("employees").
			Where(FieldYear("hire_date"),
				Between, ValueBetween{
					Low:  1990,
					High: 1993,
				}).
			OrderBy("hire_date", Asc),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT $1 OFFSET $2)": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("salary", Eq,
				QueryInstance().
					Select("DISTINCT salary").
					From("employees").
					OrderBy("salary", Desc).
					Limit(1, 1),
			),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE salary >= ALL (SELECT salary FROM employees WHERE department_id = $1) ORDER BY salary DESC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("salary", GrEqAll,
				QueryInstance().
					Select("salary").
					From("employees").
					Where("department_id", Eq, 8),
			).
			OrderBy("salary", Desc),
		"SELECT employee_id, first_name, last_name, salary FROM employees e WHERE  EXISTS (SELECT 1 FROM dependents d WHERE d.employee_id = e.employee_id) ORDER BY first_name ASC, last_name ASC": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees", "e").
			Where(FieldEmpty(""), Exists,
				QueryInstance().
					Select("1").
					From("dependents", "d").
					Where("d.employee_id", Eq, ValueField("e.employee_id")),
			).
			OrderBy("first_name", Asc).
			OrderBy("last_name", Asc),
		"SELECT employee_id, first_name, last_name, salary FROM employees WHERE department_id IN (SELECT department_id FROM departments WHERE department_name = $1 OR department_name = $2)": QueryInstance().
			Select("employee_id", "first_name", "last_name", "salary").
			From("employees").
			Where("department_id", In,
				QueryInstance().
					Select("department_id").
					From("departments").
					Where("department_name", Eq, "Marketing").
					WhereOr("department_name", Eq, "Sales"),
			),
		"SELECT inv_no AS invoice_no, amount, due_date AS 'Due date', cust_no 'Customer No' FROM invoices": QueryInstance().
			Select("inv_no AS invoice_no", "amount", "due_date AS 'Due date'", "cust_no 'Customer No'").
			From("invoices"),
		"SELECT first_name, last_name, salary * 1.1 AS new_salary FROM employees WHERE new_salary > $1": QueryInstance().
			Select("first_name", "last_name", "salary * 1.1 AS new_salary").
			From("employees").
			Where("new_salary", Greater, 5000),
		"SELECT first_name, last_name, employees.department_id, departments.department_id, department_name FROM employees INNER JOIN departments ON departments.department_id = employees.department_id WHERE employees.department_id IN ($1, $2, $3)": QueryInstance().
			Select("first_name", "last_name", "employees.department_id", "departments.department_id", "department_name").
			From("employees").
			Join(InnerJoin, "departments", Condition{
				Field: "departments.department_id",
				Opt:   Eq,
				Value: ValueField("employees.department_id"),
			}).
			Where("employees.department_id", In, []int{1, 2, 3}),
		"SELECT first_name, last_name, job_title, department_name FROM employees e INNER JOIN departments d ON d.department_id = e.department_id INNER JOIN jobs j ON j.job_id = e.job_id WHERE e.department_id IN ($1, $2, $3)": QueryInstance().
			Select("first_name", "last_name", "job_title", "department_name").
			From("employees", "e").
			Join(InnerJoin, "departments d", Condition{
				Field: "d.department_id",
				Opt:   Eq,
				Value: ValueField("e.department_id"),
			}).
			Join(InnerJoin, "jobs j", Condition{
				Field: "j.job_id",
				Opt:   Eq,
				Value: ValueField("e.job_id"),
			}).
			Where("e.department_id", In, []int{1, 2, 3}),
		"SELECT c.country_name, c.country_id, l.country_id, l.street_address, l.city FROM countries c LEFT JOIN locations l ON l.country_id = c.country_id WHERE c.country_id IN ($1, $2, $3)": QueryInstance().
			Select("c.country_name", "c.country_id", "l.country_id", "l.street_address", "l.city").
			From("countries", "c").
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("c.country_id", In, []string{"US", "UK", "CN"}),
		"SELECT country_name FROM countries c LEFT JOIN locations l ON l.country_id = c.country_id WHERE l.location_id IS NULL ORDER BY country_name ASC": QueryInstance().
			Select("country_name").
			From("countries", "c").
			Join(LeftJoin, "locations l", Condition{
				Field: "l.country_id",
				Opt:   Eq,
				Value: ValueField("c.country_id"),
			}).
			Where("l.location_id", Null, nil).
			OrderBy("country_name", Asc),
		"SELECT e.first_name || ' ' || e.last_name AS employee, m.first_name || ' ' || m.last_name AS manager FROM employees e INNER JOIN employees m ON m.employee_id = e.manager_id ORDER BY manager ASC": QueryInstance().
			Select("e.first_name || ' ' || e.last_name AS employee", "m.first_name || ' ' || m.last_name AS manager").
			From("employees", "e").
			Join(InnerJoin, "employees m", Condition{
				Field: "m.employee_id",
				Opt:   Eq,
				Value: ValueField("e.manager_id"),
			}).
			OrderBy("manager", Asc),
		"SELECT basket_name, fruit_name FROM fruits FULL OUTER JOIN baskets ON baskets.basket_id = fruits.basket_id WHERE fruit_name IS NULL": QueryInstance().
			Select("basket_name", "fruit_name").
			From("fruits").
			Join(FullOuterJoin, "baskets", Condition{
				Field: "baskets.basket_id",
				Opt:   Eq,
				Value: ValueField("fruits.basket_id"),
			}).
			Where("fruit_name", Null, nil),
	}

	for expected, query := range testCases {
		var sql string
		var args []any

		sql, args, _ = query.Sql()

		if sql != expected {
			t.Fatalf(`Query %s != %s (%v)`, sql, expected, args)
		}
	}
}

// TestRemoveLimit
func TestRemoveLimit(t *testing.T) {
	query := QueryInstance().
		Select("salary").
		From("employees").
		OrderBy("salary", Desc).
		Limit(10, 3)

	limit := query.RemoveLimit()

	if limit.Limit != 10 || limit.Offset != 3 {
		t.Fatalf(`%v`, limit)
	}
	sqlRemoveLimit := "SELECT salary FROM employees ORDER BY salary DESC"
	sqlLimit := "SELECT salary FROM employees ORDER BY salary DESC LIMIT 10 OFFSET 3"

	if sqlRemoveLimit != query.String() {
		t.Fatalf(`Query %s != %s`, query.String(), sqlRemoveLimit)
	}

	query.Limit(limit.Limit, limit.Offset)

	if sqlLimit != query.String() {
		t.Fatalf(`Query %s != %s`, query.String(), sqlLimit)
	}
}

// TestRemoveFetch
func TestRemoveFetch(t *testing.T) {
	query := QueryInstance().
		Select("salary").
		From("employees").
		OrderBy("salary", Desc).
		Fetch(3, 10)

	fetch := query.RemoveFetch()

	if fetch.Fetch != 10 || fetch.Offset != 3 {
		t.Fatalf(`%v`, fetch)
	}
	sqlRemoveLimit := "SELECT salary FROM employees ORDER BY salary DESC"
	sqlLimit := "SELECT salary FROM employees ORDER BY salary DESC OFFSET 10 ROWS FETCH NEXT 3 ROWS ONLY"

	if sqlRemoveLimit != query.String() {
		t.Fatalf(`Query %s != %s`, query.String(), sqlRemoveLimit)
	}

	query.Fetch(fetch.Fetch, fetch.Offset)

	if sqlLimit != query.String() {
		t.Fatalf(`Query %s != %s`, query.String(), sqlLimit)
	}
}
