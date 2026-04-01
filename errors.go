package lattice

import "errors"

var (
	// ErrTableNotFound indicates a table was not found.
	ErrTableNotFound = errors.New("table not found")
	// ErrColumnNotFound indicates a column was not found.
	ErrColumnNotFound = errors.New("column not found")
	// ErrCycleDetected indicates a cycle was detected in the dependency graph.
	ErrCycleDetected = errors.New("cycle detected in dependency graph")
	// ErrTypeMismatch indicates a type mismatch.
	ErrTypeMismatch = errors.New("type mismatch")
	// ErrDirtyState indicates dirty state.
	ErrDirtyState = errors.New("dirty state")
	// ErrUDFTimeout indicates a UDF timed out.
	ErrUDFTimeout = errors.New("udf timeout")
	// ErrBackfillFailed indicates a backfill operation failed.
	ErrBackfillFailed = errors.New("backfill failed")
	// ErrDuplicateTable indicates a table with the same name already exists.
	ErrDuplicateTable = errors.New("duplicate table")
	// ErrDuplicateColumn indicates a column with the same name already exists.
	ErrDuplicateColumn = errors.New("duplicate column")
	// ErrInvalidSchema indicates an invalid schema.
	ErrInvalidSchema = errors.New("invalid schema")
)
