package qb

import (
	"bytes"
	"context"
	"database/sql"
	"reflect"
	"time"
)

// TODO how to conditionally add loggedRunner
type loggedRunner struct {
	inner runner
}

func (r loggedRunner) ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	logger := GetLoggerCtx(ctx)
	if logger == nil {
		return r.inner.ExecContext(ctx, q, args...)
	}
	start := time.Now()
	result, err := r.inner.ExecContext(ctx, q, args...)
	logger(ctx, time.Since(start), "%s -- %v", q, args)
	return result, err
}

func (r loggedRunner) QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error) {
	logger := GetLoggerCtx(ctx)
	if logger == nil {
		return r.inner.QueryContext(ctx, q, args...)
	}
	start := time.Now()
	rows, err := r.inner.QueryContext(ctx, q, args...)
	logger(ctx, time.Since(start), "%s -- %v", q, args)
	return rows, err
}

func query(ctx context.Context, r runner, builder Builder, dest interface{}) (int, error) {
	buf := bytes.Buffer{}

	err := builder.Build(&buf)
	if err != nil {
		return 0, err
	}

	rows, err := loggedRunner{r}.QueryContext(ctx, buf.String(), builder.Params()...)
	if err != nil {
		return 0, err
	}

	count, err := load(rows, dest)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func exec(ctx context.Context, r runner, builder Builder) (sql.Result, error) {
	buf := bytes.Buffer{}

	err := builder.Build(&buf)
	if err != nil {
		return nil, err
	}

	result, err := loggedRunner{r}.ExecContext(ctx, buf.String(), builder.Params()...)
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
