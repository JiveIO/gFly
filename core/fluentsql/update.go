package fluentsql

import (
	"app/core/fluentsql/builder"
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

type updateData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	Table             string
	SetClauses        []setClause
	From              Sqlizer
	WhereParts        []Sqlizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

type setClause struct {
	column string
	value  interface{}
}

func (d *updateData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *updateData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}

func (d *updateData) QueryRow() RowScanner {
	if d.RunWith == nil {
		return &Row{err: RunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRower)
	if !ok {
		return &Row{err: RunnerNotQueryRunner}
	}
	return QueryRowWith(queryRower, d)
}

func (d *updateData) ToSql() (sqlOut string, args []interface{}, err error) {
	if d.Table == "" {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(d.SetClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set clause")
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

	sqlStr.WriteString("UPDATE ")
	sqlStr.WriteString(d.Table)

	sqlStr.WriteString(" SET ")
	setSqls := make([]string, len(d.SetClauses))
	for i, setClause := range d.SetClauses {
		var valSql string
		if vs, ok := setClause.value.(Sqlizer); ok {
			vsql, vargs, err := vs.ToSql()
			if err != nil {
				return "", nil, err
			}
			if _, ok := vs.(SelectBuilder); ok {
				valSql = fmt.Sprintf("(%s)", vsql)
			} else {
				valSql = vsql
			}
			args = append(args, vargs...)
		} else {
			valSql = "?"
			args = append(args, setClause.value)
		}
		setSqls[i] = fmt.Sprintf("%s = %s", setClause.column, valSql)
	}
	sqlStr.WriteString(strings.Join(setSqls, ", "))

	if d.From != nil {
		sqlStr.WriteString(" FROM ")
		args, err = appendToSql([]Sqlizer{d.From}, sqlStr, "", args)
		if err != nil {
			return
		}
	}

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

	sqlOut, err = d.PlaceholderFormat.ReplacePlaceholders(sqlStr.String())
	return
}

// Builder

// UpdateBuilder builds SQL UPDATE statements.
type UpdateBuilder builder.Builder

func init() {
	builder.Register(UpdateBuilder{}, updateData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b UpdateBuilder) PlaceholderFormat(f PlaceholderFormat) UpdateBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(UpdateBuilder)
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b UpdateBuilder) RunWith(runner BaseRunner) UpdateBuilder {
	return setRunWith(b, runner).(UpdateBuilder)
}

// Exec builds and Execs the query with the Runner set by RunWith.
func (b UpdateBuilder) Exec() (sql.Result, error) {
	data := builder.GetStruct(b).(updateData)
	return data.Exec()
}

func (b UpdateBuilder) Query() (*sql.Rows, error) {
	data := builder.GetStruct(b).(updateData)
	return data.Query()
}

func (b UpdateBuilder) QueryRow() RowScanner {
	data := builder.GetStruct(b).(updateData)
	return data.QueryRow()
}

func (b UpdateBuilder) Scan(dest ...interface{}) error {
	return b.QueryRow().Scan(dest...)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b UpdateBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(updateData)
	return data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b UpdateBuilder) MustSql() (string, []interface{}) {
	sqlStr, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sqlStr, args
}

// Prefix adds an expression to the beginning of the query
func (b UpdateBuilder) Prefix(sqlStr string, args ...interface{}) UpdateBuilder {
	return b.PrefixExpr(Expr(sqlStr, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b UpdateBuilder) PrefixExpr(expr Sqlizer) UpdateBuilder {
	return builder.Append(b, "Prefixes", expr).(UpdateBuilder)
}

// Table sets the table to be updated.
func (b UpdateBuilder) Table(table string) UpdateBuilder {
	return builder.Set(b, "Table", table).(UpdateBuilder)
}

// Set adds SET clauses to the query.
func (b UpdateBuilder) Set(column string, value interface{}) UpdateBuilder {
	return builder.Append(b, "SetClauses", setClause{column: column, value: value}).(UpdateBuilder)
}

// SetMap is a convenience method which calls .Set for each key/value pair in clauses.
func (b UpdateBuilder) SetMap(clauses map[string]interface{}) UpdateBuilder {
	keys := make([]string, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.Set(key, clauses[key])
	}
	return b
}

// From adds FROM clause to the query
// FROM is valid construct in postgresql only.
func (b UpdateBuilder) From(from string) UpdateBuilder {
	return builder.Set(b, "From", newPart(from)).(UpdateBuilder)
}

// FromSelect sets a subquery into the FROM clause of the query.
func (b UpdateBuilder) FromSelect(from SelectBuilder, alias string) UpdateBuilder {
	// Prevent misnumbered parameters in nested selects (#183).
	from = from.PlaceholderFormat(Question)
	return builder.Set(b, "From", Alias(from, alias)).(UpdateBuilder)
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b UpdateBuilder) Where(pred interface{}, args ...interface{}) UpdateBuilder {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(UpdateBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b UpdateBuilder) OrderBy(orderBys ...string) UpdateBuilder {
	return builder.Extend(b, "OrderBys", orderBys).(UpdateBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b UpdateBuilder) Limit(limit uint64) UpdateBuilder {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(UpdateBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b UpdateBuilder) Offset(offset uint64) UpdateBuilder {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(UpdateBuilder)
}

// Suffix adds an expression to the end of the query
func (b UpdateBuilder) Suffix(sqlStr string, args ...interface{}) UpdateBuilder {
	return b.SuffixExpr(Expr(sqlStr, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b UpdateBuilder) SuffixExpr(expr Sqlizer) UpdateBuilder {
	return builder.Append(b, "Suffixes", expr).(UpdateBuilder)
}
