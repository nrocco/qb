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

### Two core interfaces (`database.go`)

- **`Builder`** — implemented by all query types; has `Build(*bytes.Buffer)` (writes SQL) and `Params() []interface{}` (returns bind values).
- **`Runner`** — abstracts `*sql.DB` and `*sql.Tx`; has `ExecContext` / `QueryContext`. Both the `DB` and `Tx` wrapper structs satisfy this.

### Query builders

Each file (`select.go`, `insert.go`, `update.go`, `delete.go`) holds a query struct and fluent builder methods. All return `*QueryType` for chaining. The builder accumulates SQL fragments and params separately; rendering happens only at execution time.

- **`SelectQuery`** — `From`, `Columns`, `Where`, `Join`, `OrderBy`, `GroupBy`, `Limit`, `Offset`, `With` (CTEs). Execute with `Load(dest)` (slice/struct via reflection) or `LoadValue(dest)` (scalar).
- **`InsertQuery`** — `InTo`, `Columns`, `Values`, `Record` (struct introspection), `OrIgnore`, `OnConflict`, `Returning`. Execute with `Exec()`.
- **`UpdateQuery`** — `Table`, `Set`, `Where`, `Returning`. Execute with `Exec()`.
- **`DeleteQuery`** — `From`, `Where`. Execute with `Exec()`.

### Execution pipeline (`execute.go`)

`query()` and `exec()` are the shared helpers that build SQL, log duration via the context logger, resolve any active transaction from context, and execute against the `Runner`. `load()` uses reflection to scan rows into slices, structs, or scalars. Field mapping reads `db:` struct tags; without a tag, field names are converted CamelCase→snake_case via `structMap()` in `utils.go`.

### Transaction context (`transaction.go`)

`WitTx(ctx, tx)` stores a `*Tx` in the context. Query builders call `GetTxCtx(ctx)` at execution time and use that transaction automatically — callers don't need to thread the `Tx` through every query call explicitly.

### Null types & JSON (`types.go`, `json.go`)

`NullString`, `NullInt64`, `NullTime` wrap `sql.Null*` with proper JSON marshaling (`null` when invalid). `JSONValue` / `JSONScan` handle JSON/JSONB columns.
