package fluentsql

import (
	"testing"
)

// TestInsertTable
func TestInsertTable(t *testing.T) {
	testCases := map[string]*InsertBuilder{
		"INSERT INTO products (name, desc, category_id) VALUES ('Which Book Should I Read?', 'Which Book Should I Read?', 12)": InsertInstance().
			Insert("products", "name", "desc", "category_id").
			Row("Which Book Should I Read?", "Which Book Should I Read?", 12),
		"INSERT INTO Customers (CustomerName, City, Country) SELECT SupplierName, City, Country FROM Suppliers WHERE Country = 'Germany'": InsertInstance().
			Insert("Customers", "CustomerName", "City", "Country").
			Query(QueryInstance().
				Select("SupplierName", "City", "Country").
				From("Suppliers").
				Where("Country", Eq, "Germany"),
			),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestInsertTableArgs
func TestInsertTableArgs(t *testing.T) {
	testCases := map[string]*InsertBuilder{
		"INSERT INTO products (name, desc, category_id) VALUES ($1, $2, $3)": InsertInstance().
			Insert("products", "name", "desc", "category_id").
			Row("Which Book Should I Read?", "Which Book Should I Read?", 12),
		"INSERT INTO Customers (CustomerName, City, Country) SELECT SupplierName, City, Country FROM Suppliers WHERE Country = $1": InsertInstance().
			Insert("Customers", "CustomerName", "City", "Country").
			Query(QueryInstance().
				Select("SupplierName", "City", "Country").
				From("Suppliers").
				Where("Country", Eq, "Germany"),
			),
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
