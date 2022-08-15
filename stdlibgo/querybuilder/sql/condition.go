package sql

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Condition struct {
	stmt interface{}
}

func (c *Condition) BuildQuery() (query string, args []interface{}, err error) {
	var (
		strct   = reflect.ValueOf(c.stmt)
		typeOfT = strct.Type()
	)

	for i := 0; i < strct.NumField(); i++ {
		var (
			scopeQuery string
			scopeArg   interface{}
		)

		fieldValue := strct.Field(i).Interface()
		isValueEmpty := reflect.ValueOf(fieldValue).IsZero()

		if isValueEmpty {
			continue
		}

		fieldName := typeOfT.Field(i).Name
		fieldStrategy := FieldStrategy(fieldName)

		scopeQuery, scopeArg, err = c.buildScope(fieldStrategy, fieldValue)
		if err != nil {
			return
		}

		if query != "" {
			query += "AND "
		}

		query += scopeQuery + " "
		args = append(args, scopeArg)
	}

	if query != "" {
		query = "WHERE " + strings.Trim(query, " ")
	}

	return
}

func (c *Condition) buildScope(field FieldStrategy, value interface{}) (query string, arg interface{}, err error) {
	arg = value

	if field.IsLikeStatement() {
		query += fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field.ColumnName())
		return
	}

	if field.IsGreaterThanEqualStatement() {
		switch value.(type) {
		default:
			err = fmt.Errorf("statement %s value must be a number or date", field)
			return
		case int, int8, int32, int64, uint, uint8, uint32, uint64, float32, float64:
			query += fmt.Sprintf("%s >= ?", field.ColumnName())
			return
		}
	}

	if field.IsLessThanEqualStatement() {
		switch value.(type) {
		default:
			err = fmt.Errorf("statement %s value must be a number or date", field)
			return
		case int, int8, int32, int64, uint, uint8, uint32, uint64, float32, float64, time.Time:
			query += fmt.Sprintf("%s <= ?", field.ColumnName())
			return
		}
	}

	if field.IsNotInStatement() {
		query += fmt.Sprintf("%s NOT IN (?)", field.ColumnName())
		return
	} else if field.IsInStatement() {
		query += fmt.Sprintf("%s IN (?)", field.ColumnName())
		return
	}

	query += fmt.Sprintf("%s = ?", field.ColumnName())

	return
}
