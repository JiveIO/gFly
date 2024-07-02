package fluentsql

import (
	"testing"
)

// TestOrderBy
func TestOrderBy(t *testing.T) {
	orderByTest := new(OrderBy)

	orderByTest.Append("first_name", Asc)
	expected := "ORDER BY first_name ASC"

	if orderByTest.String() != expected {
		t.Fatalf(`Query %s != %s`, orderByTest.String(), expected)
	}
}

// TestOrderMulti
func TestOrderByTestOrderMultiDesc(t *testing.T) {
	orderByTest := new(OrderBy)

	orderByTest.Append("first_name", Asc)
	orderByTest.Append("last_name", Desc)
	expected := "ORDER BY first_name ASC, last_name DESC"

	if orderByTest.String() != expected {
		t.Fatalf(`Query %s != %s`, orderByTest.String(), expected)
	}
}
