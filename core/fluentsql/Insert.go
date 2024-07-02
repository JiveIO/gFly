package fluentsql

import (
	"fmt"
	"strings"
)

/*
INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country)
VALUES ('Cardinal', 'Tom B. Erichsen', 'Skagen 21', 'Stavanger', '4006', 'Norway');

INSERT INTO dependents (first_name, last_name, relationship, employee_id)
VALUES
	('Cameron', 'Bell', 'Child', 192),
	('Michelle', 'Bell', 'Child', 192);

INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country)
SELECT SupplierName, ContactName, Address, City, PostalCode, Country FROM Suppliers;

INSERT INTO Customers (CustomerName, City, Country)
SELECT SupplierName, City, Country FROM Suppliers WHERE Country='Germany';
*/

// Insert clause
type Insert struct {
	Table   string
	Columns []string
}

func (i *Insert) String() string {
	columnsStr := strings.Join(i.Columns, ", ")

	return fmt.Sprintf("INSERT INTO %s (%s)", i.Table, columnsStr)
}
