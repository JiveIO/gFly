package fluentsql

import (
	"fmt"
	"strings"
)

// Case clause
/*
Simple CASE expression
   CASE (2000 - YEAR(hire_date))
       WHEN 1 THEN '1 year'
       WHEN 3 THEN '3 years'
       WHEN 5 THEN '5 years'
       WHEN 10 THEN '10 years'
       WHEN 15 THEN '15 years'
       WHEN 20 THEN '20 years'
       WHEN 25 THEN '25 years'
       WHEN 30 THEN '30 years'
   END aniversary

Search CASE expression
   CASE
       WHEN salary < 3000 THEN 'Low'
       WHEN salary >= 3000 AND salary <= 5000 THEN 'Average'
       WHEN salary > 5000 THEN 'High'
   END evaluation
*/

type Case struct {
	// Exp Expression
	Exp string
	// WhenClauses expression
	WhenClauses []WhenCase
	// Case name
	Name string
}

func FieldCase(exp, name string) *Case {
	return &Case{
		Exp:  exp,
		Name: name,
	}
}

func (c *Case) When(conditions any, value string) *Case {
	c.WhenClauses = append(c.WhenClauses, WhenCase{
		Conditions: conditions,
		Value:      value,
	})

	return c
}

type WhenCase struct {
	Conditions any // string | int | []Condition
	Value      string
}

func (c *WhenCase) String() string {
	if valueConditions, ok := c.Conditions.([]Condition); ok {
		var cons []string
		for _, condition := range valueConditions {
			cons = append(cons, condition.String())
		}

		return fmt.Sprintf("WHEN %s THEN '%s'", strings.Join(cons, " AND "), c.Value)
	}

	return fmt.Sprintf("WHEN %v THEN '%s'", c.Conditions, c.Value)
}

func (c *Case) String() string {
	var whenCases []string

	for _, whenClause := range c.WhenClauses {
		whenCases = append(whenCases, whenClause.String())
	}

	return fmt.Sprintf("CASE %s %s END %s", c.Exp, strings.Join(whenCases, " "), c.Name)
}
