package fluentsql

import "testing"

// TestDelete
func TestDelete(t *testing.T) {
	testCases := map[string]Delete{
		"DELETE FROM products": {
			Table: "products",
		},
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}
