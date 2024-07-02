package fluentsql

import (
	"testing"
)

// TestFrom
func TestFrom(t *testing.T) {
	fromTest := new(From)

	fromTest.Table = "products"
	expected := "FROM products"

	if fromTest.String() != expected {
		t.Fatalf(`Query %s != %s`, fromTest.String(), expected)
	}
}

// TestFromAlias
func TestFromAlias(t *testing.T) {
	fromTest := new(From)

	fromTest.Table = "products"
	fromTest.Alias = "p"
	expected := "FROM products p"

	if fromTest.String() != expected {
		t.Fatalf(`Query %s != %s`, fromTest.String(), expected)
	}
}

// TestFromAlias
func TestFromSelect(t *testing.T) {
	fromTest := new(From)

	fromTest.Table = QueryInstance().Select("COUNT(*)").From("products").AS("counter")
	expected := "FROM (SELECT COUNT(*) FROM products) AS counter"

	if fromTest.String() != expected {
		t.Fatalf(`Query %s != %s`, fromTest.String(), expected)
	}

	fromTest.Table = QueryInstance().Select("COUNT(*)").From("products")
	fromTest.Alias = "counter"
	expected = "FROM (SELECT COUNT(*) FROM products) counter"

	if fromTest.String() != expected {
		t.Fatalf(`Query %s != %s`, fromTest.String(), expected)
	}
}
