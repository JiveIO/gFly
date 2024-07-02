package fluentmodel

import (
	"app/core/fluentsql"
	"errors"
	"reflect"
)

// Delete perform delete data for table via model type Struct, *Struct
//
// -------- Delete via Model --------
//
//	var user User
//	err = db.First(&user)
//	err = db.Delete(user)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// -------- Delete via ID --------
//
//	err = db.Delete(User{}, 157)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	-------- Delete via List ID --------
//
//	err = db.Delete(User{}, []int{154, 155, 156})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// -------- Delete via Where condition --------
//
//	err = db.Where("Id", fluentsql.Eq, 153).Delete(&User{})
//	if err != nil {
//		log.Fatal(err)
//	}
func (db *DBModel) Delete(model any, args ...any) (err error) {
	// Delete by raw sql
	if db.raw.sqlStr != "" {
		err = db.execRaw(db.raw.sqlStr, db.raw.args)

		if err != nil {
			panic(err)
		}

		// Reset fluent model builder
		db.reset()

		return
	}

	if len(args) > 0 {
		var opt fluentsql.WhereOpt
		var argument = args[0]
		argumentType := reflect.TypeOf(argument)

		if argumentType.Kind() == reflect.Slice {
			opt = fluentsql.In
		} else {
			opt = fluentsql.Eq
		}

		db.wherePrimaryCondition = fluentsql.Condition{
			Field: nil,
			Opt:   opt,
			Value: argument,
			AndOr: fluentsql.And,
		}
	}

	var table *Table
	var primaryKey any
	var hasCondition = false

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		return err
	}

	// Get a primary key
	if len(table.Primaries) > 0 {
		primaryKey = table.Primaries[0].Name

		if table.Values[primaryKey.(string)] != nil {
			db.wherePrimaryCondition = fluentsql.Condition{
				Field: nil,
				Opt:   fluentsql.Eq,
				Value: table.Values[primaryKey.(string)],
				AndOr: fluentsql.And,
			}
		}
	}

	// Create delete instance
	deleteBuilder := fluentsql.DeleteInstance().
		Delete(table.Name)

	// Build WHERE condition with specific primary value
	if db.wherePrimaryCondition.Value != nil && primaryKey != nil {
		deleteBuilder.Where(primaryKey, db.wherePrimaryCondition.Opt, db.wherePrimaryCondition.Value)
		hasCondition = true
	}

	// Build WHERE condition from a condition list
	for _, condition := range db.whereStatement.Conditions {
		switch {
		// Sub-conditions
		case len(condition.Group) > 0:
			// Append conditions from a group to query builder
			deleteBuilder.WhereGroup(func(whereBuilder fluentsql.WhereBuilder) *fluentsql.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)

				return &whereBuilder
			})
			hasCondition = true
		case condition.AndOr == fluentsql.And:
			// Add Where AND condition
			deleteBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == fluentsql.Or:
			// Add Where OR condition
			deleteBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		default:
			hasCondition = false
		}
	}

	if !hasCondition {
		panic(errors.New("missing WHERE condition for deleting operator"))
	}

	err = db.delete(deleteBuilder)

	// Reset fluent model builder
	db.reset()

	return
}
