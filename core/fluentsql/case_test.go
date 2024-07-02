package fluentsql

import "testing"

// TestCaseSimple
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
*/
func TestCaseSimple(t *testing.T) {
	caseTest := new(Case)

	caseTest.Exp = "(2000 - YEAR(hire_date))"
	caseTest.When(1, "1 year")
	caseTest.When(3, "3 years")
	caseTest.When(5, "5 years")
	caseTest.When(10, "10 years")
	caseTest.When(15, "15 years")
	caseTest.When(20, "20 years")
	caseTest.When(25, "25 years")
	caseTest.When(30, "30 years")
	caseTest.Name = "aniversary"

	expected := "CASE (2000 - YEAR(hire_date)) WHEN 1 THEN '1 year' WHEN 3 THEN '3 years' WHEN 5 THEN '5 years' WHEN 10 THEN '10 years' WHEN 15 THEN '15 years' WHEN 20 THEN '20 years' WHEN 25 THEN '25 years' WHEN 30 THEN '30 years' END aniversary"

	if caseTest.String() != expected {
		t.Fatalf(`Query %s != %s`, caseTest.String(), expected)
	}
}

// TestCaseSimple
/*
Search CASE expression
   CASE
       WHEN salary < 3000 THEN 'Low'
       WHEN salary >= 3000 AND salary <= 5000 THEN 'Average'
       WHEN salary > 5000 THEN 'High'
   END evaluation
*/
func TestCaseSearch(t *testing.T) {
	caseTest := new(Case)

	var conditionsLow []Condition
	conditionsLow = append(conditionsLow, Condition{
		Field: "salary",
		Opt:   Lesser,
		Value: 3000,
	})

	var conditionsAverage []Condition
	conditionsAverage = append(conditionsAverage, Condition{
		Field: "salary",
		Opt:   GrEq,
		Value: 3000,
	}, Condition{
		Field: "salary",
		Opt:   LeEq,
		Value: 5000,
	})

	var conditionsHigh []Condition
	conditionsHigh = append(conditionsHigh, Condition{
		Field: "salary",
		Opt:   Greater,
		Value: 5000,
	})

	caseTest.Exp = ""
	caseTest.When(conditionsLow, "Low")
	caseTest.When(conditionsAverage, "Average")
	caseTest.When(conditionsHigh, "High")
	caseTest.Name = "evaluation"

	expected := "CASE  WHEN salary < 3000 THEN 'Low' WHEN salary >= 3000 AND salary <= 5000 THEN 'Average' WHEN salary > 5000 THEN 'High' END evaluation"

	if caseTest.String() != expected {
		t.Fatalf(`Query %s != %s`, caseTest.String(), expected)
	}
}
