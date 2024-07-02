package fluentsql

import (
	"testing"
)

// TestGroupBy
func TestGroupBy(t *testing.T) {
	groupByTest := new(GroupBy)

	groupByTest.Append("first_name")

	expected := "GROUP BY first_name"

	if groupByTest.String() != expected {
		t.Fatalf(`Query %s != %s`, groupByTest.String(), expected)
	}
}

// TestGroupByColumns
func TestGroupByColumns(t *testing.T) {
	groupByTest := new(GroupBy)

	groupByTest.Append("first_name")
	groupByTest.Append("last_name")

	expected := "GROUP BY first_name, last_name"

	if groupByTest.String() != expected {
		t.Fatalf(`Query %s != %s`, groupByTest.String(), expected)
	}
}
