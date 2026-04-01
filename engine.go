package lattice

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
)

// Engine is the top-level orchestrator for Lattice.
type Engine struct {
	store       Store
	blobStore   BlobStore
	vectorStore VectorStore
	logger      *slog.Logger
	tables      map[string]*Table
}

// New creates a new Engine with the given options.
func New(opts ...Option) *Engine {
	e := &Engine{
		logger: slog.Default(),
		tables: make(map[string]*Table),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// CreateTable creates a new table from a schema.
func (e *Engine) CreateTable(ctx context.Context, schema TableSchema) (*Table, error) {
	if err := schema.Validate(); err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	if _, exists := e.tables[schema.Name]; exists {
		return nil, fmt.Errorf("creating table %q: %w", schema.Name, ErrDuplicateTable)
	}
	if e.store != nil {
		if err := e.store.CreateTable(ctx, schema); err != nil {
			return nil, fmt.Errorf("creating table %q in store: %w", schema.Name, err)
		}
	}
	t := &Table{schema: schema, engine: e}
	e.tables[schema.Name] = t
	return t, nil
}

// Table returns a table by name.
func (e *Engine) Table(name string) (*Table, error) {
	t, ok := e.tables[name]
	if !ok {
		return nil, fmt.Errorf("getting table %q: %w", name, ErrTableNotFound)
	}
	return t, nil
}

// DropTable removes a table.
func (e *Engine) DropTable(ctx context.Context, name string) error {
	if _, ok := e.tables[name]; !ok {
		return fmt.Errorf("dropping table %q: %w", name, ErrTableNotFound)
	}
	if e.store != nil {
		if err := e.store.DropTable(ctx, name); err != nil {
			return fmt.Errorf("dropping table %q from store: %w", name, err)
		}
	}
	delete(e.tables, name)
	return nil
}

// Tables returns all table names.
func (e *Engine) Tables() []string {
	names := make([]string, 0, len(e.tables))
	for name := range e.tables {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
