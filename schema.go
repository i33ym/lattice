package lattice

import "fmt"

// ColumnSpec defines a column in a table schema.
type ColumnSpec struct {
	Name     string
	Type     Type
	Computed *ComputedSpec
}

// TableSchema defines the schema for a table.
type TableSchema struct {
	Name    string
	Columns []ColumnSpec
}

// Validate checks the schema for errors such as duplicate names, missing types, and invalid computed references.
func (s TableSchema) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("validating schema: %w: table name is empty", ErrInvalidSchema)
	}
	if len(s.Columns) == 0 {
		return fmt.Errorf("validating schema %q: %w: no columns defined", s.Name, ErrInvalidSchema)
	}

	seen := make(map[string]bool, len(s.Columns))
	for _, col := range s.Columns {
		if seen[col.Name] {
			return fmt.Errorf("validating schema %q: %w: %q", s.Name, ErrDuplicateColumn, col.Name)
		}
		seen[col.Name] = true

		if col.Type == nil && col.Computed == nil {
			return fmt.Errorf("validating schema %q column %q: %w: nil type on non-computed column", s.Name, col.Name, ErrInvalidSchema)
		}
	}

	for _, col := range s.Columns {
		if col.Computed == nil {
			continue
		}
		for _, input := range col.Computed.Inputs {
			if !seen[input] {
				return fmt.Errorf("validating schema %q computed column %q: %w: references non-existent column %q", s.Name, col.Name, ErrColumnNotFound, input)
			}
		}
	}

	return nil
}
