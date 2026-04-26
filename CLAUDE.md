# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build       # lint + test (default target)
make lint        # staticcheck, golint, go vet
make test        # go test -v -short ./...
make coverage    # HTML coverage report in coverage/
```

Run a single test:
```bash
go test -v -run TestSelectQuery ./...
```

## Architecture

**qb** is a SQLite query builder for Go (`github.com/nrocco/qb`). It uses `modernc.org/sqlite` (pure Go — no CGo).

### Separation of concerns

Query builders and execution are fully decoupled:

- **Builders** (`select.go`, `insert.go`, `update.go`, `delete.go`) are pure value objects. They implement the `Builder` interface (`Build(*bytes.Buffer) error` + `Params() []interface{}`), accumulate SQL state via fluent method chains, and have no reference to any database connection.
- **`DB` and `Tx`** (`database.go`, `transaction.go`) own all execution. They provide factory methods that return bare builders, and execution methods that take a `Builder`:

```go
// Build — no DB reference needed
q := db.Select().From("notes").Where("id = ?", 1).Limit(10)

// Execute
count, err := db.Load(ctx, q, &notes)       // into slice/struct
err  = db.LoadValue(ctx, q, &scalar)         // single value
result, err := db.Exec(ctx, q)               // insert/update/delete
id, _ := result.LastInsertId()
```

### DB vs Tx execution

`DB.Exec/Load/LoadValue` auto-sniff a `*Tx` from context via `GetTxCtx(ctx)` and use it if present. `Tx.Exec/Load/LoadValue` always use the transaction's runner directly. To pass a transaction via context: `WithTx(ctx, tx)`.

### Shared WHERE clause (`where.go`)

`SelectQuery`, `UpdateQuery`, and `DeleteQuery` all embed `whereClause`, which provides `addWhere()` and `writeWhere()`. Each builder's `Where()` method delegates to `addWhere` and returns its own concrete type for fluent chaining.

### Execution pipeline (`execute.go`)

`query()` and `exec()` are the shared helpers. Both wrap the incoming runner with `loggedRunner`, which checks `GetLoggerCtx(ctx)` — if `nil`, the call passes through with zero overhead; otherwise it times and logs the query. Inject a logger via `WithLogger(ctx, fn)`.

### Struct scanning (`execute.go`, `utils.go`)

`load()` uses reflection to scan rows into slices, structs, or scalars. Field mapping reads `db:` struct tags; without a tag, field names are converted CamelCase→snake_case via `structMap()`. `InsertQuery.Record()` uses the same `structMap` to populate values from a struct.
