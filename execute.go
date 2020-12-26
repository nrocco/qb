package qb

import (
	"bytes"
	"context"
	"database/sql"
	"reflect"
)

func query(ctx context.Context, runner runner, builder Builder, dest interface{}) (int, error) {
	buf := bytes.Buffer{}

	err := builder.Build(&buf)
	if err != nil {
		return 0, err
	}

	query, params := buf.String(), builder.Params()

	logger := GetLogCtx(ctx)
	logger("%s -- %v", query, params)

	rows, err := runner.QueryContext(ctx, query, params...)
	if err != nil {
		return 0, err
	}

	count, err := load(rows, dest)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func exec(ctx context.Context, runner runner, builder Builder) (sql.Result, error) {
	buf := bytes.Buffer{}

	err := builder.Build(&buf)
	if err != nil {
		return nil, err
	}

	query, params := buf.String(), builder.Params()

	logger := GetLogCtx(ctx)
	logger("%s -- %v", query, params)

	result, err := runner.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func load(rows *sql.Rows, value interface{}) (int, error) {
	defer rows.Close()

	column, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return 0, ErrInvalidPointer
	}

	v = v.Elem()
	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8
	count := 0

	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(v.Type().Elem()).Elem()
		} else {
			elem = v
		}

		ptr, err := findPtr(column, elem)
		if err != nil {
			return 0, err
		}

		err = rows.Scan(ptr...)
		if err != nil {
			return 0, err
		}

		count++

		if isSlice {
			v.Set(reflect.Append(v, elem))
		} else {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

type dummyScanner struct{}

func (dummyScanner) Scan(interface{}) error {
	return nil
}

var (
	dummyDest   sql.Scanner = dummyScanner{}
	typeScanner             = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
)

func findPtr(column []string, value reflect.Value) ([]interface{}, error) {
	if value.Addr().Type().Implements(typeScanner) {
		return []interface{}{value.Addr().Interface()}, nil
	}
	switch value.Kind() {
	case reflect.Struct:
		var ptr []interface{}
		m := structMap(value)
		for _, key := range column {
			if val, ok := m[key]; ok {
				ptr = append(ptr, val.Addr().Interface())
			} else {
				ptr = append(ptr, dummyDest)
			}
		}
		return ptr, nil
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return findPtr(column, value.Elem())
	}
	return []interface{}{value.Addr().Interface()}, nil
}
