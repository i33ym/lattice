package lattice

import (
	"context"
	"errors"
	"fmt"
)

// Table represents a table in the engine.
type Table struct {
	schema TableSchema
	engine *Engine
}

// Name returns the table name.
func (t *Table) Name() string { return t.schema.Name }

// Schema returns the table schema.
func (t *Table) Schema() TableSchema { return t.schema }

// Insert inserts rows into the table.
func (t *Table) Insert(ctx context.Context, rows []Row) error {
	if t.engine.store == nil {
		return errors.New("inserting rows: no store configured")
	}
	if err := t.engine.store.Insert(ctx, t.schema.Name, rows); err != nil {
		return fmt.Errorf("inserting rows into %q: %w", t.schema.Name, err)
	}
	return nil
}

// Select queries the table with an optional filter expression.
func (t *Table) Select(ctx context.Context, filter Expr) (RowIterator, error) {
	if t.engine.store == nil {
		return nil, errors.New("querying rows: no store configured")
	}
	iter, err := t.engine.store.Query(ctx, t.schema.Name, filter)
	if err != nil {
		return nil, fmt.Errorf("querying %q: %w", t.schema.Name, err)
	}
	return iter, nil
}

// AddColumn adds a new column to the table schema.
func (t *Table) AddColumn(ctx context.Context, spec ColumnSpec) error {
	for _, col := range t.schema.Columns {
		if col.Name == spec.Name {
			return fmt.Errorf("adding column %q to %q: %w", spec.Name, t.schema.Name, ErrDuplicateColumn)
		}
	}
	t.schema.Columns = append(t.schema.Columns, spec)
	if err := t.schema.Validate(); err != nil {
		t.schema.Columns = t.schema.Columns[:len(t.schema.Columns)-1]
		return fmt.Errorf("adding column %q to %q: %w", spec.Name, t.schema.Name, err)
	}
	return nil
}

// AlterColumn modifies an existing column in the table schema.
func (t *Table) AlterColumn(ctx context.Context, name string, spec ColumnSpec) error {
	for i, col := range t.schema.Columns {
		if col.Name == name {
			t.schema.Columns[i] = spec
			return nil
		}
	}
	return fmt.Errorf("altering column %q in %q: %w", name, t.schema.Name, ErrColumnNotFound)
}
