package fluentsql

import (
	"testing"
)

// TestDeleteTable
func TestDeleteTable(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' AND city = 'Frankfurt' AND customer_id = 1": DeleteInstance().
			Delete("customers").
			Where("contact_name", Eq, "Alfred Schmidt").
			Where("city", Eq, "Frankfurt").
			Where("customer_id", Eq, 1),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestDeleteTableArgs
func TestDeleteTableArgs(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = $1 AND city = $2 AND customer_id = $3": DeleteInstance().
			Delete("customers").
			Where("contact_name", Eq, "Alfred Schmidt").
			Where("city", Eq, "Frankfurt").
			Where("customer_id", Eq, 1),
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
