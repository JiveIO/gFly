package fluentsql

import (
	"testing"
)

// TestUpdateBasic
func TestUpdateBasic(t *testing.T) {
	testCases := map[string]*UpdateBuilder{
		"UPDATE Customers SET ContactName = 'Alfred Schmidt', City = 'Frankfurt' WHERE CustomerID = 1": UpdateInstance().
			Update("Customers").
			Set("ContactName", "Alfred Schmidt").
			Set("City", "Frankfurt").
			Where("CustomerID", Eq, 1),
		"UPDATE Customers SET ContactName = 'Alfred Schmidt', City = Location WHERE CustomerID = 5": UpdateInstance().
			Update("Customers").
			Set("ContactName", "Alfred Schmidt").
			Set("City", ValueField("Location")).
			Where("CustomerID", Eq, 5),
		"UPDATE CUSTOMERS SET ADDRESS = 'Pune', SALARY = 1000 WHERE NAME = 'Ramesh'": UpdateInstance().
			Update("CUSTOMERS").
			Set("ADDRESS", "Pune").
			Set("SALARY", 1000.00).
			Where("NAME", Eq, "Ramesh"),
		"UPDATE dependents SET last_name = (SELECT last_name FROM employees WHERE employee_id = dependents.employee_id)": UpdateInstance().
			Update("dependents").
			Set("last_name", QueryInstance().
				Select("last_name").
				From("employees").
				Where("employee_id", Eq, ValueField("dependents.employee_id")),
			),
		"UPDATE summary s SET (sum_x, sum_y, avg_x, avg_y) = (SELECT sum(x), sum(y), avg(x), avg(y) FROM data d WHERE d.group_id = s.group_id)": UpdateInstance().
			Update("summary", "s").
			Set([]string{"sum_x", "sum_y", "avg_x", "avg_y"}, QueryInstance().
				Select("sum(x)", "sum(y)", "avg(x)", "avg(y)").
				From("data", "d").
				Where("d.group_id", Eq, ValueField("s.group_id")),
			),
		"UPDATE summary SET (sum_x, sum_y, avg_x) = (1, 'One', 34.5)": UpdateInstance().
			Update("summary").
			Set([]string{"sum_x", "sum_y", "avg_x"}, []any{1, "One", 34.5}),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestUpdateArgs
func TestUpdateArgs(t *testing.T) {
	testCases := map[string]*UpdateBuilder{
		"UPDATE Customers SET ContactName = $1, City = $2 WHERE CustomerID = $3": UpdateInstance().
			Update("Customers").
			Set("ContactName", "Alfred Schmidt").
			Set("City", "Frankfurt").
			Where("CustomerID", Eq, 1),
		"UPDATE Customers SET ContactName = $1, City = Location WHERE CustomerID = $2": UpdateInstance().
			Update("Customers").
			Set("ContactName", "Alfred Schmidt").
			Set("City", ValueField("Location")).
			Where("CustomerID", Eq, 5),
		"UPDATE CUSTOMERS SET ADDRESS = $1, SALARY = $2 WHERE NAME = $3": UpdateInstance().
			Update("CUSTOMERS").
			Set("ADDRESS", "Pune").
			Set("SALARY", 1000.00).
			Where("NAME", Eq, "Ramesh"),
		"UPDATE dependents SET last_name = (SELECT last_name FROM employees WHERE employee_id = dependents.employee_id)": UpdateInstance().
			Update("dependents").
			Set("last_name", QueryInstance().
				Select("last_name").
				From("employees").
				Where("employee_id", Eq, ValueField("dependents.employee_id")),
			),
		"UPDATE summary s SET (sum_x, sum_y, avg_x, avg_y) = (SELECT sum(x), sum(y), avg(x), avg(y) FROM data d WHERE d.group_id = s.group_id)": UpdateInstance().
			Update("summary", "s").
			Set([]string{"sum_x", "sum_y", "avg_x", "avg_y"}, QueryInstance().
				Select("sum(x)", "sum(y)", "avg(x)", "avg(y)").
				From("data", "d").
				Where("d.group_id", Eq, ValueField("s.group_id")),
			),
		"UPDATE summary SET (sum_x, sum_y, avg_x) = ($1, $2, $3)": UpdateInstance().
			Update("summary").
			Set([]string{"sum_x", "sum_y", "avg_x"}, []any{1, "One", 34.5}),
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
