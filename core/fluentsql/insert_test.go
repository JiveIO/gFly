package fluentsql

import (
	"testing"
)

// TestInsert
func TestInsert(t *testing.T) {
	testCases := map[string]Insert{
		"INSERT INTO products (first_name, last_name, category_id)": {
			Table:   "products",
			Columns: []string{"first_name", "last_name", "category_id"},
		},
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}

// TestInsertRow
func TestInsertRow(t *testing.T) {
	testCases := map[string]InsertRow{
		"(first_name, 'last_name', 12, 92.3)": {
			Values: []any{ValueField("first_name"), "last_name", 12, 92.3},
		},
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}

// TestInsertRows
func TestInsertRows(t *testing.T) {
	insertRow := InsertRows{}
	insertRow.Append(ValueField("first_name"), "last_name", 12, 92.3)
	insertRow.Append("first_name", "last_name", "12", 35.3)

	testCases := map[string]InsertRows{
		"VALUES (first_name, 'last_name', 12, 92.3), ('first_name', 'last_name', '12', 35.3)": insertRow,
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}
