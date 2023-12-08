package fluentsql

import (
	"app/core/fluentsql/builder"
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type deleteData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	From              string
	WhereParts        []Sqlizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

func (d *deleteData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *deleteData) ToSql() (sqlStrOut string, args []interface{}, err error) {
	if d.From == "" {
		err = fmt.Errorf("delete statements must specify a From table")
		return
	}

	sqlStr := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSql(d.Prefixes, sqlStr, " ", args)
		if err != nil {
			return
		}

		sqlStr.WriteString(" ")
	}

	sqlStr.WriteString("DELETE FROM ")
	sqlStr.WriteString(d.From)

	if len(d.WhereParts) > 0 {
		sqlStr.WriteString(" WHERE ")
		args, err = appendToSql(d.WhereParts, sqlStr, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderBys) > 0 {
		sqlStr.WriteString(" ORDER BY ")
		sqlStr.WriteString(strings.Join(d.OrderBys, ", "))
	}

	if len(d.Limit) > 0 {
		sqlStr.WriteString(" LIMIT ")
		sqlStr.WriteString(d.Limit)
	}

	if len(d.Offset) > 0 {
		sqlStr.WriteString(" OFFSET ")
		sqlStr.WriteString(d.Offset)
	}

	if len(d.Suffixes) > 0 {
		sqlStr.WriteString(" ")
		args, err = appendToSql(d.Suffixes, sqlStr, " ", args)
		if err != nil {
			return
		}
	}

	sqlStrOut, err = d.PlaceholderFormat.ReplacePlaceholders(sqlStr.String())
	return
}

// Builder

// DeleteBuilder builds SQL DELETE statements.
type DeleteBuilder builder.Builder

func init() {
	builder.Register(DeleteBuilder{}, deleteData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b DeleteBuilder) PlaceholderFormat(f PlaceholderFormat) DeleteBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(DeleteBuilder)
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b DeleteBuilder) RunWith(runner BaseRunner) DeleteBuilder {
	return setRunWith(b, runner).(DeleteBuilder)
}

// Exec builds and Execs the query with the Runner set by RunWith.
func (b DeleteBuilder) Exec() (sql.Result, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.Exec()
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b DeleteBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b DeleteBuilder) MustSql() (string, []interface{}) {
	sqlStr, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sqlStr, args
}

// Prefix adds an expression to the beginning of the query
func (b DeleteBuilder) Prefix(sqlStr string, args ...interface{}) DeleteBuilder {
	return b.PrefixExpr(Expr(sqlStr, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b DeleteBuilder) PrefixExpr(expr Sqlizer) DeleteBuilder {
	return builder.Append(b, "Prefixes", expr).(DeleteBuilder)
}

// From sets the table to be deleted from.
func (b DeleteBuilder) From(from string) DeleteBuilder {
	return builder.Set(b, "From", from).(DeleteBuilder)
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b DeleteBuilder) Where(pred interface{}, args ...interface{}) DeleteBuilder {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(DeleteBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b DeleteBuilder) OrderBy(orderBys ...string) DeleteBuilder {
	return builder.Extend(b, "OrderBys", orderBys).(DeleteBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b DeleteBuilder) Limit(limit uint64) DeleteBuilder {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(DeleteBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b DeleteBuilder) Offset(offset uint64) DeleteBuilder {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(DeleteBuilder)
}

// Suffix adds an expression to the end of the query
func (b DeleteBuilder) Suffix(sqlStr string, args ...interface{}) DeleteBuilder {
	return b.SuffixExpr(Expr(sqlStr, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b DeleteBuilder) SuffixExpr(expr Sqlizer) DeleteBuilder {
	return builder.Append(b, "Suffixes", expr).(DeleteBuilder)
}

func (b DeleteBuilder) Query() (*sql.Rows, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.Query()
}

func (d *deleteData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}
