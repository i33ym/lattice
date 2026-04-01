package lattice

import "context"

// Store is the interface for relational data storage.
type Store interface {
	CreateTable(ctx context.Context, schema TableSchema) error
	DropTable(ctx context.Context, name string) error
	Insert(ctx context.Context, table string, rows []Row) error
	Update(ctx context.Context, table string, id string, cols map[string]any) error
	Delete(ctx context.Context, table string, id string) error
	Query(ctx context.Context, table string, filter Expr) (RowIterator, error)
	Schema(ctx context.Context, table string) (TableSchema, error)
	Get(ctx context.Context, table string, id string) (Row, error)
}
