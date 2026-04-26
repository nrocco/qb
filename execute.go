package qb

import (
	"bytes"
	"context"
	"database/sql"
	"reflect"
	"sync"
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

var bufPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

func query(ctx context.Context, r runner, builder Builder, dest interface{}) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := builder.Build(buf); err != nil {
		return 0, err
	}

	rows, err := loggedRunner{r}.QueryContext(ctx, buf.String(), builder.Params()...)
	if err != nil {
		return 0, err
	}

	return load(rows, dest)
}

func exec(ctx context.Context, r runner, builder Builder) (sql.Result, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := builder.Build(buf); err != nil {
		return nil, err
	}

	return loggedRunner{r}.ExecContext(ctx, buf.String(), builder.Params()...)
}

// scanPlan is a precomputed mapping from result column positions to struct field paths.
// A nil entry means no matching field (use dummyDest).
type scanPlan [][]int

func newScanPlan(columns []string, t reflect.Type) scanPlan {
	nameToPath := make(map[string][]int)
	buildFieldPaths(t, nil, nameToPath)

	plan := make(scanPlan, len(columns))
	for i, col := range columns {
		plan[i] = nameToPath[col]
	}
	return plan
}

func buildFieldPaths(t reflect.Type, prefix []int, out map[string][]int) {
	if reflect.PointerTo(t).Implements(typeValuer) || t.Implements(typeValuer) {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}
		tag := field.Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = camelCaseToSnakeCase(field.Name)
		}

		path := make([]int, len(prefix)+1)
		copy(path, prefix)
		path[len(prefix)] = i

		if _, exists := out[tag]; !exists {
			out[tag] = path
		}

		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			buildFieldPaths(ft, path, out)
		}
	}
}

func load(rows *sql.Rows, value interface{}) (int, error) {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return 0, ErrInvalidPointer
	}
	v = v.Elem()

	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8

	var elemType reflect.Type
	if isSlice {
		elemType = v.Type().Elem()
	} else {
		elemType = v.Type()
	}

	// Dereference pointer element type for plan building
	baseType := elemType
	if baseType.Kind() == reflect.Ptr {
		baseType = baseType.Elem()
	}

	// Precompute column→field plan for structs; fall back to findPtr for scanners/scalars
	var plan scanPlan
	if baseType.Kind() == reflect.Struct && !reflect.PointerTo(baseType).Implements(typeScanner) {
		plan = newScanPlan(columns, baseType)
	}

	// Pre-allocate the ptrs slice once and reuse across rows
	ptrs := make([]interface{}, len(columns))
	count := 0

	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(elemType).Elem()
		} else {
			elem = v
		}

		if plan != nil {
			target := elem
			if target.Kind() == reflect.Ptr {
				if target.IsNil() {
					target.Set(reflect.New(target.Type().Elem()))
				}
				target = target.Elem()
			}
			for i, path := range plan {
				if path == nil {
					ptrs[i] = dummyDest
					continue
				}
				f := target
				for _, idx := range path {
					if f.Kind() == reflect.Ptr {
						if f.IsNil() {
							f.Set(reflect.New(f.Type().Elem()))
						}
						f = f.Elem()
					}
					f = f.Field(idx)
				}
				ptrs[i] = f.Addr().Interface()
			}
			if err = rows.Scan(ptrs...); err != nil {
				return 0, err
			}
		} else {
			p, err := findPtr(columns, elem)
			if err != nil {
				return 0, err
			}
			if err = rows.Scan(p...); err != nil {
				return 0, err
			}
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
