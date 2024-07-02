package fluentmodel

import (
	"app/core/fluentsql"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"    // load driver for Mysql
	_ "github.com/jackc/pgx/v5/stdlib"    // load driver for PostgreSQL
	_ "github.com/joho/godotenv/autoload" // load .env file
)

// ===========================================================================================================
// 											Structure & Interface
// ===========================================================================================================

// Init Fluent Model's actions.
func init() {
	initDB()

	// More initial functions.
}

type iDatabase interface {
	connect() (*sqlx.DB, error)
}

var (
	// Define database connection settings.
	maxConn         = GetEnv("DB_MAX_CONNECTION", 0)          // the default is 0 (unlimited)
	maxIdleConn     = GetEnv("DB_MAX_IDLE_CONNECTION", 2)     // default is 2
	maxLifetimeConn = GetEnv("DB_MAX_LIFETIME_CONNECTION", 0) // default is 0, connections are reused forever
)

// connectDB connect to Database.
func connectDB(connURL, driver string) (*sqlx.DB, error) {
	// Define database connection.
	dbConnection, err := sqlx.Connect(driver, connURL)
	if err != nil {
		return nil, err
	}

	// Set database connection settings:
	dbConnection.SetMaxOpenConns(maxConn)
	dbConnection.SetMaxIdleConns(maxIdleConn)
	dbConnection.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := dbConnection.Ping(); err != nil {
		defer func(db *sqlx.DB) {
			_ = db.Close()
		}(dbConnection)
		return nil, err
	}

	return dbConnection, nil
}

// ===========================================================================================================
// 											PostgreSQL Connection
// ===========================================================================================================

// postgreSQL a implement of interface iDatabase for PostgreSQL
type postgreSQL struct{}

// Connect perform DB connection to PostgreSQL database.
func (db *postgreSQL) connect() (*sqlx.DB, error) {
	// Build PostgreSQL connection URL.
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	return connectDB(connURL, "pgx")
}

// ===========================================================================================================
// 											MySQL Connection
// ===========================================================================================================

// mySQL a implement of interface iDatabase for MySQL
type mySQL struct{}

// Connect perform DB connection to Mysql database.
func (db *mySQL) connect() (*sqlx.DB, error) {
	// Build Mysql connection URL.
	connURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return connectDB(connURL, "mysql")
}

// ===========================================================================================================
// 												Database
// ===========================================================================================================

// defaultDB a singleton database instance
var dbInstance = &DB{}

// DBInstanceTx returns db transaction instance to handle CRUD at somewhere.
func dbInstanceTx() *sqlx.Tx {
	return dbInstance.MustBegin()
}

// DB the database
type DB struct {
	*sqlx.DB // Embed sqlx DB.
}

// Connect func for opening database connection.
func initDB() {
	var err error
	var database iDatabase

	// Get DB_TYPE value from .env file.
	driverType := GetEnv("DB_DRIVER", "mysql")
	dbType := fluentsql.SQLite

	// Define a new Database connection with a right DB type.
	switch driverType {
	case "postgresql":
		database = &postgreSQL{}
		dbType = fluentsql.PostgreSQL
	case "mysql":
		database = &mySQL{}
		dbType = fluentsql.MySQL
	default:
		panic("Unknown database driver")
	}

	// Initial DBType
	fluentsql.SetDBType(dbType)

	dbInstance.DB, err = database.connect()
	if err != nil {
		panic(err)
	}
}

// ===========================================================================================================
// 												DB Model
// ===========================================================================================================

// Raw struct
type Raw struct {
	sqlStr string
	args   []any
}

type DBModel struct {
	tx *sqlx.Tx

	model any // Model struct
	raw   Raw // Raw struct

	selectStatement       fluentsql.Select    // Select columns
	omitsSelectStatement  fluentsql.Select    // Omit columns
	whereStatement        fluentsql.Where     // Where conditions
	wherePrimaryCondition fluentsql.Condition // Where primary key condition
	joinStatement         fluentsql.Join
	groupByStatement      fluentsql.GroupBy
	havingStatement       fluentsql.Having // A version of Where
	orderByStatement      fluentsql.OrderBy
	limitStatement        fluentsql.Limit
	fetchStatement        fluentsql.Fetch // A version of Limit
}

func Instance() *DBModel {
	return &DBModel{
		tx:    nil,
		model: nil,
	}
}

// Reset DB model's builders after everytime perform the DB query
func (db *DBModel) reset() *DBModel {
	db.model = nil
	db.raw.sqlStr = ""
	db.selectStatement.Columns = []any{}
	db.omitsSelectStatement.Columns = []any{}
	db.whereStatement.Conditions = []fluentsql.Condition{}
	db.wherePrimaryCondition.Value = nil
	db.joinStatement.Items = []fluentsql.JoinItem{}
	db.groupByStatement.Items = []string{}
	db.havingStatement.Conditions = []fluentsql.Condition{}
	db.orderByStatement.Items = []fluentsql.SortItem{}
	db.limitStatement.Limit = 0
	db.fetchStatement.Fetch = 0

	return db
}

// ===========================================================================================================
// 											FluentSQL + SQLX integration
// ===========================================================================================================

// get perform getting single data row by QueryBuilder
func (db *DBModel) get(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.getRaw(sqlStr, args, model)
}

// get perform getting single data row by QueryBuilder
func (db *DBModel) getRaw(sqlStr string, args []any, model any) (err error) {
	if GetEnv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	if db.tx != nil {
		err = db.tx.Get(model, sqlStr, args...)
	} else {
		err = dbInstance.Get(model, sqlStr, args...)
	}

	return
}

// query performs query list data row by QueryBuilder
func (db *DBModel) query(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.queryRaw(sqlStr, args, model)
}

// queryRaw performs query list data row by sqlStr and arguments
func (db *DBModel) queryRaw(sqlStr string, args []any, model any) (err error) {
	if GetEnv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	if db.tx != nil {
		err = db.tx.Select(model, sqlStr, args...)
	} else {
		err = dbInstance.Select(model, sqlStr, args...)
	}

	return
}

// add performs adding new data by InsertBuilder
func (db *DBModel) add(q *fluentsql.InsertBuilder, primaryColumn string) (id any, err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.addRaw(sqlStr, args, primaryColumn)
}

// addRaw performs adding new data by sqlStr and arguments
func (db *DBModel) addRaw(sqlStr string, args []any, primaryColumn string) (id any, err error) {
	if GetEnv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	// Data persistence
	if fluentsql.DBType() == fluentsql.PostgreSQL {
		if primaryColumn != "" {
			sqlStr += " RETURNING " + primaryColumn

			if GetEnv("DB_DEBUG", false) {
				log.Printf("Chagned SQL> %s", sqlStr)
			}
		}

		if db.tx != nil {
			err = db.tx.QueryRow(sqlStr, args...).Scan(&id)
		} else {
			err = dbInstance.QueryRow(sqlStr, args...).Scan(&id)
		}
	} else if fluentsql.DBType() == fluentsql.MySQL {
		var result sql.Result
		if db.tx != nil {
			result, err = db.tx.Exec(sqlStr, args...)
			if err != nil {
				return nil, err
			}
		} else {
			result, err = dbInstance.Exec(sqlStr, args...)
			if err != nil {
				return nil, err
			}
		}

		id, err = result.LastInsertId()
	}

	return
}

// update performs updating data by UpdateBuilder
func (db *DBModel) update(q *fluentsql.UpdateBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// delete performs deleting data by DeleteBuilder
func (db *DBModel) delete(q *fluentsql.DeleteBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// execRaw performs updating and deleting data by DeleteBuilder
func (db *DBModel) execRaw(sqlStr string, args []any) (err error) {
	if GetEnv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	// Data persistence
	if db.tx != nil {
		_, err = db.tx.Exec(sqlStr, args...)
	} else {
		_, err = dbInstance.Exec(sqlStr, args...)
	}

	return
}

// Count get total rows
func (db *DBModel) count(q *fluentsql.QueryBuilder, total *int) error {
	var fetch fluentsql.Fetch
	var limit fluentsql.Limit

	// Build SQL without pagination
	fetch = q.RemoveFetch()
	limit = q.RemoveLimit()

	// Create COUNT query
	sqlBuilderCount := fluentsql.QueryInstance().
		Select("COUNT(*) AS total").
		From(q, "_result_out_")

	err := db.get(sqlBuilderCount, total)
	if err != nil {
		return err
	}

	// Reset pagination
	q.Limit(limit.Limit, limit.Offset)
	q.Fetch(fetch.Offset, fetch.Fetch)

	return nil
}

// ===========================================================================================================
// 												DB Model operators
// ===========================================================================================================

// Raw build query from raw SQL
func (db *DBModel) Raw(sqlStr string, args ...any) *DBModel {
	db.raw.sqlStr = sqlStr
	db.raw.args = args

	return db
}

// Select List of columns
func (db *DBModel) Select(columns ...any) *DBModel {
	db.selectStatement.Columns = columns

	return db
}

// Omit exclude some columns
func (db *DBModel) Omit(columns ...any) *DBModel {
	db.omitsSelectStatement.Columns = columns

	return db
}

// Model set specific model for builder
func (db *DBModel) Model(model any) *DBModel {
	db.model = model

	return db
}

// Where add where condition
func (db *DBModel) Where(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// WhereOr add where condition
func (db *DBModel) WhereOr(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.Or,
	})

	return db
}

// WhereGroup combine multi where conditions into a group.
func (db *DBModel) WhereGroup(groupCondition fluentsql.FnWhereBuilder) *DBModel {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	cond := fluentsql.Condition{
		Group: whereBuilder.Conditions(),
	}

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, cond)

	return db
}

// When checking TRUE to build Where condition.
func (db *DBModel) When(condition bool, groupCondition fluentsql.FnWhereBuilder) *DBModel {
	if !condition {
		return db
	}

	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, whereBuilder.Conditions()...)

	return db
}

// Join builder
func (db *DBModel) Join(join fluentsql.JoinType, table string, condition fluentsql.Condition) *DBModel {
	db.joinStatement.Append(fluentsql.JoinItem{
		Join:      join,
		Table:     table,
		Condition: condition,
	})

	return db
}

// Having builder
func (db *DBModel) Having(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.havingStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// GroupBy fields in a query
func (db *DBModel) GroupBy(fields ...string) *DBModel {
	db.groupByStatement.Append(fields...)

	return db
}

// OrderBy builder
func (db *DBModel) OrderBy(field string, dir fluentsql.OrderByDir) *DBModel {
	db.orderByStatement.Append(field, dir)

	return db
}

// Limit builder
func (db *DBModel) Limit(limit, offset int) *DBModel {
	db.limitStatement.Limit = limit
	db.limitStatement.Offset = offset

	return db
}

// RemoveLimit builder
func (db *DBModel) RemoveLimit() fluentsql.Limit {
	var _limitStatement fluentsql.Limit

	_limitStatement.Limit = db.limitStatement.Limit
	_limitStatement.Offset = db.limitStatement.Offset

	db.limitStatement.Limit = 0
	db.limitStatement.Offset = 0

	return _limitStatement
}

// Fetch builder
func (db *DBModel) Fetch(offset, fetch int) *DBModel {
	db.fetchStatement.Offset = offset
	db.fetchStatement.Fetch = fetch

	return db
}

// RemoveFetch builder
func (db *DBModel) RemoveFetch() fluentsql.Fetch {
	var _fetchStatement fluentsql.Fetch

	_fetchStatement.Offset = db.fetchStatement.Offset
	_fetchStatement.Fetch = db.fetchStatement.Fetch

	db.fetchStatement.Offset = 0
	db.fetchStatement.Fetch = 0

	return _fetchStatement
}

// whereFromModel Build and append WHERE clause from specific model's data off table.
func (tbl *Table) whereFromModel(queryBuilder *fluentsql.QueryBuilder) {
	if tbl.HasData {
		for _, column := range tbl.Columns {
			// Prevent some meta, relational, and default (Zero) value of column
			if column.isNotData() || column.IsZero {
				continue
			}

			// Append query conditions
			queryBuilder.Where(column.Name, fluentsql.Eq, tbl.Values[column.Name])
		}
	}
}

// ===========================================================================================================
// 												DB Transaction
// ===========================================================================================================

// Begin new transaction
func (db *DBModel) Begin() *DBModel {
	db.tx = dbInstanceTx()

	return db
}

// Rollback transaction
func (db *DBModel) Rollback() error {
	if db.tx != nil {
		return db.tx.Rollback()
	}

	return nil
}

// Commit transaction
func (db *DBModel) Commit() error {
	if db.tx != nil {
		return db.tx.Commit()
	}

	return nil
}
