package dirty

import "sync"

// CellID identifies a single cell by its row and column.
type CellID struct {
	RowID  string
	Column string
}

// Tracker maintains a set of dirty cells that need recomputation.
type Tracker struct {
	cells map[CellID]bool
	mu    sync.RWMutex
}

// New creates a Tracker with no dirty cells.
func New() *Tracker {
	return &Tracker{
		cells: make(map[CellID]bool),
	}
}

// Mark flags a single cell as dirty.
func (t *Tracker) Mark(rowID, column string) {
	t.mu.Lock()
	t.cells[CellID{RowID: rowID, Column: column}] = true
	t.mu.Unlock()
}

// MarkMany flags multiple columns for a given row as dirty.
func (t *Tracker) MarkMany(rowID string, columns []string) {
	t.mu.Lock()
	for _, col := range columns {
		t.cells[CellID{RowID: rowID, Column: col}] = true
	}
	t.mu.Unlock()
}

// Clear removes the dirty flag from a single cell.
func (t *Tracker) Clear(rowID, column string) {
	t.mu.Lock()
	delete(t.cells, CellID{RowID: rowID, Column: column})
	t.mu.Unlock()
}

// IsDirty reports whether a specific cell is marked dirty.
func (t *Tracker) IsDirty(rowID, column string) bool {
	t.mu.RLock()
	dirty := t.cells[CellID{RowID: rowID, Column: column}]
	t.mu.RUnlock()
	return dirty
}

// DirtyCells returns all currently dirty cells.
func (t *Tracker) DirtyCells() []CellID {
	t.mu.RLock()
	result := make([]CellID, 0, len(t.cells))
	for cell := range t.cells {
		result = append(result, cell)
	}
	t.mu.RUnlock()
	return result
}

// DirtyColumns returns the dirty column names for a specific row.
func (t *Tracker) DirtyColumns(rowID string) []string {
	t.mu.RLock()
	var columns []string
	for cell := range t.cells {
		if cell.RowID == rowID {
			columns = append(columns, cell.Column)
		}
	}
	t.mu.RUnlock()
	return columns
}

// DirtyRows returns the dirty row IDs for a specific column.
func (t *Tracker) DirtyRows(column string) []string {
	t.mu.RLock()
	var rows []string
	for cell := range t.cells {
		if cell.Column == column {
			rows = append(rows, cell.RowID)
		}
	}
	t.mu.RUnlock()
	return rows
}

// Count returns the number of dirty cells.
func (t *Tracker) Count() int {
	t.mu.RLock()
	n := len(t.cells)
	t.mu.RUnlock()
	return n
}

// Reset clears all dirty flags.
func (t *Tracker) Reset() {
	t.mu.Lock()
	t.cells = make(map[CellID]bool)
	t.mu.Unlock()
}
