/*
Convert from model struct to table structure.
*/

package fluentmodel

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ===========================================================================================================
// 												Struct and data
// ===========================================================================================================

const (
	MODEL     = "model"   // Tag `model`
	TABLE     = "table"   // Table name
	TYPE      = "type"    // Column types
	REFERENCE = "ref"     // Column reference
	CASCADE   = "cascade" // Column cascade DELETE, UPDATE
	RELATION  = "rel"     // Column relationship
	NAME      = "name"    // Column name
)

// MetaData name
type MetaData string

// Table structure
type Table struct {
	Name      string
	Columns   []Column
	Primaries []Column
	Values    map[string]any
	Relation  []*Table
	HasData   bool
}

// Column structure
type Column struct {
	Key      string
	Name     string
	Primary  bool
	Types    string
	Ref      string // Reference id to table
	Relation string // Relation to table
	IsZero   bool   // Keep Zero value of type
	HasValue bool
}

// isNotData Does the column is column data of table
func (c *Column) isNotData() bool {
	return !c.HasValue || c.Relation != "" || c.Ref != ""
}

func NewTable() *Table {
	tbl := new(Table)
	tbl.Values = make(map[string]any)

	return tbl
}

func ModelData(model any) (*Table, error) {
	// Make sure that the input is a struct having any other type,
	// especially a pointer to a struct, might result in panic
	typ := reflect.TypeOf(model)
	value := reflect.ValueOf(model)

	// Create a table struct
	tbl := NewTable()

	// Set table data from model struct
	tbl.Name = toSnakeCase(typ.Name())

	// Port data from model to table struct
	if typ.Kind() == reflect.Struct {
		return processModel(typ, value, tbl), nil
	} else if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		return processModel(typ.Elem(), value, tbl), nil
	}

	return nil, errors.New("input param should be a struct")
}

// ===========================================================================================================
// 												Process methods
// ===========================================================================================================

func processModel(typ reflect.Type, value reflect.Value, tbl *Table) *Table {
	// Pointer case
	value = reflect.Indirect(value)

	// Go rough via field number
	for i := 0; i < typ.NumField(); i++ {
		var col Column
		// Get type field
		typeField := typ.Field(i)

		// Get value field
		var valueField reflect.Value
		// FIX panic: reflect: Field index out of range
		if value.NumField() == typ.NumField() {
			valueField = value.Field(i)
		} else {
			valueField = reflect.ValueOf(nil)
		}

		// Extract attributes in model via tab name MODEL
		attr := readTags(typeField.Tag.Get(MODEL))

		// Check primary column
		isPrimaryColumn := isPrimary(attr[TYPE])

		// Try to process special column type MetaData for Model
		if typeField.Type == reflect.TypeOf(MetaData("")) {
			if slice, tableOk := attr[TABLE]; tableOk && len(slice) > 0 {
				tbl.Name = slice[0]
			}

			continue
		}

		// Set column name
		if slice, nameOk := attr[NAME]; nameOk && len(slice) > 0 {
			col.Name = slice[0]
		} else {
			// Default name
			col.Name = toSnakeCase(typeField.Name)
		}
		col.Key = typeField.Name

		// Most functions and methods never return an invalid Value.
		validValue := valueField.IsValid()

		// Prevent primary column get Zero value
		validValueType := (isPrimaryColumn && valueField.CanInt() && !valueField.IsZero()) || !isPrimaryColumn

		// Value of column
		if validValue && validValueType {
			tbl.Values[col.Name] = valueField.Interface()
			tbl.HasData = true
			col.HasValue = true
			col.IsZero = valueField.IsZero()
		} else {
			tbl.Values[col.Name] = nil
			col.HasValue = false
			col.IsZero = true
		}

		// Process column types
		if slice, typeOk := attr[TYPE]; typeOk {
			col.Types = getTypes(slice)
		}

		// TODO Need to review ORM for REFERENCE, CASCADE, RELATION

		// Process references
		if slice, refOk := attr[REFERENCE]; refOk {
			col.Ref = getReferences(slice[0], col.Name)
		}

		// Process cascade
		if slice, casOk := attr[CASCADE]; casOk {
			col.Ref += getCascade(slice)
		}

		// Process relations
		if slice, relOk := attr[RELATION]; relOk && validValue {
			col.Relation = toSnakeCase(slice[0])

			if typeField.Type.Kind() == reflect.Slice {
				for n := 0; n < valueField.Len(); n++ {
					elemVal := valueField.Index(i)
					_tbl, _ := ModelData(elemVal.Interface())

					tbl.Relation = append(tbl.Relation, _tbl)
				}
			} else {
				_tbl, _ := ModelData(valueField.Interface())

				tbl.Relation = append(tbl.Relation, _tbl)
			}
		}

		// Column is a primary key
		col.Primary = isPrimaryColumn
		if col.Primary {
			tbl.Primaries = append(tbl.Primaries, col)
		}

		tbl.Columns = append(tbl.Columns, col)
	}

	return tbl
}

func readTags(tags string) map[string][]string {
	if tags == "" {
		return map[string][]string{TYPE: {"BOOLEAN"}}
	}

	tags = strings.ReplaceAll(tags, " ", "")
	attributes := strings.Split(tags, ";")
	var vals = make(map[string][]string)

	for i := 0; i < len(attributes); i++ {
		pre := strings.SplitN(attributes[i], ":", 2)
		vals[pre[0]] = strings.Split(pre[1], ",")
	}

	return vals
}

func getTypes(slice []string) (out string) {
	for i := 0; i < len(slice); i++ {
		var t string
		switch slice[i] {
		case "primary":
			t = "PRIMARY KEY"
		default:
			t = strings.ToUpper(slice[i])
		}
		out += t + " "
	}

	return
}

func isPrimary(slice []string) bool {
	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case "primary":
			return true
		}
	}

	return false
}

func getReferences(item, colName string) string {
	tName := toSnakeCase(item)
	refColum := strings.SplitN(colName, "_", 2)

	return fmt.Sprintf("REFERENCES %s (%s) ", tName, refColum[1])
}

func getCascade(slice []string) (out string) {
	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case "delete":
			out += "ON DELETE CASCADE "
		default:
			out += "ON UPDATE CASCADE "
		}
	}

	return
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
