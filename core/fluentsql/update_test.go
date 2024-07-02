package fluentsql

import (
	"testing"
)

// TestUpdate
func TestUpdate(t *testing.T) {
	testCases := map[string]Update{
		"UPDATE products": {
			Table: "products",
		},
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}

// TestUpdateItem
func TestUpdateItem(t *testing.T) {
	testCases := map[string]UpdateItem{
		"first_name = e.first_name": {
			Field: "first_name",
			Value: ValueField("e.first_name"),
		},
		"first_name = (SELECT first_name FROM users)": {
			Field: "first_name",
			Value: QueryInstance().
				Select("first_name").
				From("users"),
		},
	}

	for expected, item := range testCases {
		if item.String() != expected {
			t.Fatalf(`Query %s != %s`, item.String(), expected)
		}
	}
}

// TestUpdateSet
func TestUpdateSet(t *testing.T) {
	testCases := map[string]UpdateSet{
		"SET first_name = e.first_name, email = e.email": {
			Items: []UpdateItem{
				{
					Field: "first_name",
					Value: ValueField("e.first_name"),
				},
				{
					Field: "email",
					Value: ValueField("e.email"),
				},
			},
		},
		"SET first_name = e.first_name, first_name = (SELECT first_name FROM users)": {
			Items: []UpdateItem{
				{
					Field: "first_name",
					Value: ValueField("e.first_name"),
				},
				{
					Field: "first_name",
					Value: QueryInstance().
						Select("first_name").
						From("users"),
				},
			},
		},
	}

	for expected, item := range testCases {
		if item.String() != expected {
			t.Fatalf(`Query %s != %s`, item.String(), expected)
		}
	}
}
