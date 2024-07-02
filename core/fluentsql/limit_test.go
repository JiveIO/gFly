package fluentsql

import (
	"testing"
)

// TestLimit
func TestLimit(t *testing.T) {
	limitTest := new(Limit)

	limitTest.Limit = 10
	limitTest.Offset = 0
	expected := "LIMIT 10 OFFSET 0"

	if limitTest.String() != expected {
		t.Fatalf(`Query %s != %s`, limitTest.String(), expected)
	}
}

// TestFetch
func TestFetch(t *testing.T) {
	limitTest := new(Fetch)

	limitTest.Fetch = 10
	limitTest.Offset = 0
	expected := "OFFSET 0 ROWS FETCH NEXT 10 ROWS ONLY"

	if limitTest.String() != expected {
		t.Fatalf(`Query %s != %s`, limitTest.String(), expected)
	}
}
