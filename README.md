# Lattice

A reactive, pluggable multimodal data engine. Define computed columns as functions of other columns, and Lattice automatically maintains a dependency graph, tracks dirty state, and recomputes downstream values when inputs change.

Think of it as a smart spreadsheet for ML/AI pipelines — insert raw data, define transformations declaratively, and the system keeps everything in sync.

```go
engine := lattice.New(lattice.WithStore(pgStore))

table, _ := engine.CreateTable(ctx, lattice.TableSchema{
    Name: "podcasts",
    Columns: []lattice.ColumnSpec{
        {Name: "audio", Type: lattice.AudioType},
        {Name: "transcript", Computed: lattice.Computed(transcribeUDF, "audio")},
        {Name: "summary", Computed: lattice.Computed(summarizeUDF, "transcript")},
        {Name: "embedding", Computed: lattice.Computed(embedUDF, "summary")},
    },
})

table.Insert(ctx, []lattice.Row{
    {"audio": lattice.BlobFromURL("https://example.com/ep1.mp3")},
})
// transcript, summary, and embedding compute automatically
```

## Install

```bash
go get github.com/i33ym/lattice@latest
```

## Core Concepts

**DAG Engine** — Columns form a directed acyclic graph. When a source column changes, the engine walks the graph and recomputes only what's needed, in the correct order.

**Computed Columns** — Define a column as a function of other columns via a UDF. Multi-column inputs are supported: `Computed(udf, "col1", "col2")`.

**Multimodal Type System** — First-class support for `StringType`, `IntType`, `FloatType`, `BoolType`, `ImageType`, `AudioType`, `VideoType`, `VectorType(n)`, `ListOf(T)`, `JSONType`, and `TimestampType`. Register custom types at runtime with `TypeRegistry`.

**Dirty Tracking** — State is tracked per `(row_id, column_name)` cell. Only dirty cells get recomputed — no wasted work.

**UDFs** — User-Defined Functions run as external gRPC services. Write them in any language. Lattice makes zero assumptions about your models, frameworks, or algorithms.

**Backfill** — When you add or alter a computed column, Lattice backfills existing rows automatically.

## Architecture

```text
lattice/                  Root package — pure interfaces and types. Zero external deps.
├── internal/
│   ├── dag/              DAG construction, topological sort, cycle detection
│   ├── dirty/            Dirty cell tracker — (row, column) granularity
│   ├── dispatch/         Work dispatch — default goroutine pool
│   ├── eval/             UDF evaluator with retry and batching
│   ├── backfill/         Backfill planner and executor
│   └── query/            Query builder and planner
├── latticetest/          Mock implementations and test helpers
├── udf/                  UDF contract (protobuf + gRPC client)
├── store/                Storage adapters (separate go.mod each)
├── blobstore/            Blob storage adapters
├── vectorstore/          Vector storage adapters
├── dispatch/             Dispatch adapters (NATS, etc.)
├── auth/                 Auth middleware adapters
├── server/               gRPC + HTTP API layer
├── cmd/latticed/         Daemon entry point
├── cmd/lattice/          CLI entry point
├── sdk/                  Python and TypeScript SDKs
└── playground/           Web UI
```

## Adapters

| Adapter    | Type        | Package                | Status  |
| ---------- | ----------- | ---------------------- | ------- |
| PostgreSQL | Store       | `store/postgres`       | Planned |
| SQLite     | Store       | `store/sqlite`         | Planned |
| MinIO      | BlobStore   | `blobstore/minio`      | Planned |
| Local FS   | BlobStore   | `blobstore/fs`         | Planned |
| ChromaDB   | VectorStore | `vectorstore/chroma`   | Planned |
| pgvector   | VectorStore | `vectorstore/pgvector` | Planned |
| NATS       | Dispatcher  | `dispatch/nats`        | Planned |
| OpenFGA    | Auth        | `auth/openfga`         | Planned |
| Ory        | Auth        | `auth/ory`             | Planned |

Each adapter is a separate Go module with its own `go.mod`. The root module has **zero external dependencies**.

## Interfaces

Lattice is built around three core storage interfaces. Implement them to plug in any backend:

```go
type Store interface {
    CreateTable(ctx context.Context, schema TableSchema) error
    Insert(ctx context.Context, table string, rows []Row) error
    Query(ctx context.Context, table string, filter Expr) (RowIterator, error)
    // + Update, Delete, Get, DropTable, Schema
}

type BlobStore interface {
    Put(ctx context.Context, key string, r io.Reader, meta BlobMeta) error
    Get(ctx context.Context, key string) (io.ReadCloser, BlobMeta, error)
    // + Delete, Exists
}

type VectorStore interface {
    Upsert(ctx context.Context, collection string, records []VectorRecord) error
    Search(ctx context.Context, collection string, vector []float32, k int) ([]VectorResult, error)
    // + CreateCollection, DropCollection, Delete
}
```

## Development

```bash
make test          # run unit tests
make test/race     # run with race detector
make test/cover    # run with coverage report
make lint          # run golangci-lint
make audit         # vet + staticcheck + tests
make build/all     # build daemon and CLI
```

## Query Expressions

Filter rows using composable expressions:

```go
results, _ := table.Select(ctx, lattice.NewAnd(
    lattice.NewEq("status", "published"),
    lattice.NewGT("score", 0.8),
))

// Vector similarity search
results, _ = table.Select(ctx, lattice.NewVectorSearch("embedding", queryVec, 10))
```

## Playground

Open `playground/index.html` in a browser to explore the interactive web UI — includes a landing page, live playground, DAG visualizer, and admin console.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Write tests for your changes
4. Ensure `make audit` passes
5. Submit a pull request

## License

[MIT License](LICENSE)
