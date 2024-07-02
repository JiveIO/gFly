package fluentsql

import (
	"testing"
)

// TestSingleColumn
func TestSingleColumn(t *testing.T) {
	selectTest := new(Select)

	selectTest.Columns = []any{"user_name"}
	expected := "SELECT user_name"

	if selectTest.String() != expected {
		t.Fatalf(`Query %s != %s`, selectTest.String(), expected)
	}
}

// TestMultiColumns
func TestMultiColumns(t *testing.T) {
	selectTest := new(Select)

	selectTest.Columns = []any{"user_name", "email"}
	expected := "SELECT user_name, email"

	if selectTest.String() != expected {
		t.Fatalf(`Query %s != %s`, selectTest.String(), expected)
	}
}

// TestColumnFromSelect
func TestColumnFromSelect(t *testing.T) {
	selectTest := new(Select)

	selectTest.Columns = []any{"user_name", QueryInstance().Select("COUNT(*)").From("products").AS("counter")}
	expected := "SELECT user_name, (SELECT COUNT(*) FROM products) AS counter"

	if selectTest.String() != expected {
		t.Fatalf(`Query %s != %s`, selectTest.String(), expected)
	}

	selectTest.Columns = []any{"user_name", QueryInstance().Select("COUNT(*)").From("products")}
	expected = "SELECT user_name, (SELECT COUNT(*) FROM products)"

	if selectTest.String() != expected {
		t.Fatalf(`Query %s != %s`, selectTest.String(), expected)
	}
}
