package fluentsql

import (
	"app/core/fluentsql/builder"
	"bytes"
	"errors"
)

func init() {
	builder.Register(CaseBuilder{}, caseData{})
}

// sqlizerBuffer is a helper that allows to write many Sqlizers one by one
// without constant checks for errors that may come from Sqlizer
type sqlizerBuffer struct {
	bytes.Buffer
	args []interface{}
	err  error
}

// WriteSql converts Sqlizer to SQL strings and writes it to buffer
func (b *sqlizerBuffer) WriteSql(item Sqlizer) {
	if b.err != nil {
		return
	}

	var str string
	var args []interface{}
	str, args, b.err = nestedToSql(item)

	if b.err != nil {
		return
	}

	_, _ = b.WriteString(str)
	_ = b.WriteByte(' ')
	b.args = append(b.args, args...)
}

func (b *sqlizerBuffer) ToSql() (string, []interface{}, error) {
	return b.String(), b.args, b.err
}

// whenPart is a helper structure to describe SQLs "WHEN ... THEN ..." expression
type whenPart struct {
	when Sqlizer
	then Sqlizer
}

func newWhenPart(when, then interface{}) whenPart {
	return whenPart{newPart(when), newPart(then)}
}

// caseData holds all the data required to build a CASE SQL construct
type caseData struct {
	What      Sqlizer
	WhenParts []whenPart
	Else      Sqlizer
}

// ToSql implements Sqlizer
func (d *caseData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.WhenParts) == 0 {
		err = errors.New("case expression must contain at lease one WHEN clause")

		return
	}

	sql := sqlizerBuffer{}

	_, _ = sql.WriteString("CASE ")
	if d.What != nil {
		sql.WriteSql(d.What)
	}

	for _, p := range d.WhenParts {
		_, _ = sql.WriteString("WHEN ")
		sql.WriteSql(p.when)
		_, _ = sql.WriteString("THEN ")
		sql.WriteSql(p.then)
	}

	if d.Else != nil {
		_, _ = sql.WriteString("ELSE ")
		sql.WriteSql(d.Else)
	}

	_, _ = sql.WriteString("END")

	return sql.ToSql()
}

// CaseBuilder builds SQL CASE construct which could be used as parts of queries.
type CaseBuilder builder.Builder

// ToSql builds the query into a SQL string and bound args.
func (b CaseBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(caseData)
	return data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b CaseBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// what sets optional value for CASE construct "CASE [value] ..."
func (b CaseBuilder) what(expr interface{}) CaseBuilder {
	return builder.Set(b, "What", newPart(expr)).(CaseBuilder)
}

// When adds "WHEN ... THEN ..." part to CASE construct
func (b CaseBuilder) When(when, then interface{}) CaseBuilder {
	// performance hint: replace slice of WhenPart with just slice of parts
	// where even indices of the slice belong to "when"s and odd indices belong to "then"s
	return builder.Append(b, "WhenParts", newWhenPart(when, then)).(CaseBuilder)
}

// Else What sets optional "ELSE ..." part for CASE construct
func (b CaseBuilder) Else(expr interface{}) CaseBuilder {
	return builder.Set(b, "Else", newPart(expr)).(CaseBuilder)
}
