package lattice

// Row represents a single row of data as a map of column name to value.
type Row map[string]any

// RowIterator allows iterating over query results.
type RowIterator interface {
	Next() bool
	Row() Row
	Err() error
	Close() error
}

// CollectRows drains a RowIterator into a slice.
func CollectRows(iter RowIterator) ([]Row, error) {
	var rows []Row
	for iter.Next() {
		rows = append(rows, iter.Row())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return rows, nil
}
