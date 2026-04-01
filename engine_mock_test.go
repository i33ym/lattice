package lattice

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestEngineWithMockStoreFullFlow(t *testing.T) {
	t.Parallel()

	t.Run("InsertFlow/withMockStore/delegatesToStore", func(t *testing.T) {
		t.Parallel()

		var insertCalled atomic.Bool
		var capturedTable string
		var capturedRowCount int

		mock := &stubStoreCapture{
			insertFunc: func(_ context.Context, table string, rows []Row) error {
				insertCalled.Store(true)
				capturedTable = table
				capturedRowCount = len(rows)
				return nil
			},
		}

		e := New(WithStore(mock))
		ctx := context.Background()

		schema := TableSchema{
			Name: "events",
			Columns: []ColumnSpec{
				{Name: "id", Type: StringType},
				{Name: "payload", Type: JSONType},
			},
		}

		tbl, err := e.CreateTable(ctx, schema)
		if err != nil {
			t.Fatalf("CreateTable() unexpected error: %v", err)
		}

		rows := []Row{
			{"id": "evt-1", "payload": "{}"},
			{"id": "evt-2", "payload": "{}"},
		}
		if err := tbl.Insert(ctx, rows); err != nil {
			t.Fatalf("Insert() unexpected error: %v", err)
		}

		if !insertCalled.Load() {
			t.Fatalf("Store.Insert was not called")
		}
		if capturedTable != "events" {
			t.Fatalf("Store.Insert table = %q, want %q", capturedTable, "events")
		}
		if capturedRowCount != 2 {
			t.Fatalf("Store.Insert row count = %d, want 2", capturedRowCount)
		}
	})

	t.Run("InsertFlow/storeReturnsError/propagatesError", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("disk full")
		mock := &stubStoreCapture{
			insertFunc: func(context.Context, string, []Row) error {
				return wantErr
			},
		}

		e := New(WithStore(mock))
		ctx := context.Background()

		schema := TableSchema{
			Name: "logs",
			Columns: []ColumnSpec{
				{Name: "id", Type: StringType},
			},
		}

		tbl, err := e.CreateTable(ctx, schema)
		if err != nil {
			t.Fatalf("CreateTable() unexpected error: %v", err)
		}

		err = tbl.Insert(ctx, []Row{{"id": "1"}})
		if err == nil {
			t.Fatalf("Insert() expected error, got nil")
		}
		if !errors.Is(err, wantErr) {
			t.Fatalf("Insert() error = %v, want wrapping %v", err, wantErr)
		}
	})

	t.Run("InsertFlow/noStoreConfigured/returnsError", func(t *testing.T) {
		t.Parallel()

		e := New()
		ctx := context.Background()

		schema := TableSchema{
			Name: "nostore",
			Columns: []ColumnSpec{
				{Name: "id", Type: StringType},
			},
		}

		tbl, err := e.CreateTable(ctx, schema)
		if err != nil {
			t.Fatalf("CreateTable() unexpected error: %v", err)
		}

		err = tbl.Insert(ctx, []Row{{"id": "1"}})
		if err == nil {
			t.Fatalf("Insert() with no store expected error, got nil")
		}
	})

	t.Run("CreateTable/storeCreateFails/propagatesError", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("store create failed")
		mock := &stubStoreCapture{
			createTableFunc: func(context.Context, TableSchema) error {
				return wantErr
			},
		}

		e := New(WithStore(mock))
		_, err := e.CreateTable(context.Background(), TableSchema{
			Name:    "fail",
			Columns: []ColumnSpec{{Name: "id", Type: StringType}},
		})
		if err == nil {
			t.Fatalf("CreateTable() expected error, got nil")
		}
		if !errors.Is(err, wantErr) {
			t.Fatalf("CreateTable() error = %v, want wrapping %v", err, wantErr)
		}
	})
}

type stubStoreCapture struct {
	stubStore
	createTableFunc func(context.Context, TableSchema) error
	insertFunc      func(context.Context, string, []Row) error
}

func (s *stubStoreCapture) CreateTable(ctx context.Context, schema TableSchema) error {
	if s.createTableFunc != nil {
		return s.createTableFunc(ctx, schema)
	}
	return nil
}

func (s *stubStoreCapture) Insert(ctx context.Context, table string, rows []Row) error {
	if s.insertFunc != nil {
		return s.insertFunc(ctx, table, rows)
	}
	return nil
}
