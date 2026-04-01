package lattice

import (
	"context"
	"errors"
	"testing"
)

func validSchema(name string) TableSchema {
	return TableSchema{
		Name: name,
		Columns: []ColumnSpec{
			{Name: "id", Type: StringType},
			{Name: "value", Type: IntType},
		},
	}
}

func TestNewEngine(t *testing.T) {
	t.Parallel()

	t.Run("New/defaults/returnsNonNil", func(t *testing.T) {
		t.Parallel()
		e := New()
		if e == nil {
			t.Fatalf("New() returned nil")
		}
		if len(e.Tables()) != 0 {
			t.Fatalf("New() Tables() = %d, want 0", len(e.Tables()))
		}
	})

	t.Run("New/withStore/setsStore", func(t *testing.T) {
		t.Parallel()
		store := &stubStore{}
		e := New(WithStore(store))
		if e == nil {
			t.Fatalf("New(WithStore) returned nil")
		}
	})
}

func TestEngineCreateTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		schemas []TableSchema
		wantErr error
	}{
		{
			"CreateTable/validSchema/succeeds",
			[]TableSchema{validSchema("test")},
			nil,
		},
		{
			"CreateTable/duplicateName/returnsErrDuplicateTable",
			[]TableSchema{validSchema("dup"), validSchema("dup")},
			ErrDuplicateTable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := New()
			var lastErr error
			for _, s := range tt.schemas {
				_, lastErr = e.CreateTable(context.Background(), s)
			}
			if tt.wantErr == nil {
				if lastErr != nil {
					t.Fatalf("CreateTable() unexpected error: %v", lastErr)
				}
				return
			}
			if lastErr == nil {
				t.Fatalf("CreateTable() expected error wrapping %v, got nil", tt.wantErr)
			}
			if !errors.Is(lastErr, tt.wantErr) {
				t.Fatalf("CreateTable() error = %v, want errors.Is %v", lastErr, tt.wantErr)
			}
		})
	}
}

func TestEngineTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		createName string
		lookupName string
		wantErr    error
	}{
		{"Table/existingTable/returnsTable", "items", "items", nil},
		{"Table/nonExistingTable/returnsErrTableNotFound", "items", "missing", ErrTableNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := New()
			_, err := e.CreateTable(context.Background(), validSchema(tt.createName))
			if err != nil {
				t.Fatalf("CreateTable() unexpected error: %v", err)
			}
			tbl, err := e.Table(tt.lookupName)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Table(%q) error = %v, want errors.Is %v", tt.lookupName, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Table(%q) unexpected error: %v", tt.lookupName, err)
			}
			if tbl.Name() != tt.lookupName {
				t.Fatalf("Table(%q).Name() = %q", tt.lookupName, tbl.Name())
			}
		})
	}
}

func TestEngineDropTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		create  string
		drop    string
		wantErr error
	}{
		{"DropTable/existingTable/succeeds", "droppable", "droppable", nil},
		{"DropTable/nonExistingTable/returnsErrTableNotFound", "other", "missing", ErrTableNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := New()
			_, err := e.CreateTable(context.Background(), validSchema(tt.create))
			if err != nil {
				t.Fatalf("CreateTable() unexpected error: %v", err)
			}
			err = e.DropTable(context.Background(), tt.drop)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("DropTable(%q) error = %v, want errors.Is %v", tt.drop, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("DropTable(%q) unexpected error: %v", tt.drop, err)
			}
			_, err = e.Table(tt.drop)
			if !errors.Is(err, ErrTableNotFound) {
				t.Fatalf("Table(%q) after drop: error = %v, want ErrTableNotFound", tt.drop, err)
			}
		})
	}
}

func TestEngineTables(t *testing.T) {
	t.Parallel()

	t.Run("Tables/multipleCreated/returnsSortedNames", func(t *testing.T) {
		t.Parallel()
		e := New()
		ctx := context.Background()
		for _, name := range []string{"zeta", "alpha", "middle"} {
			if _, err := e.CreateTable(ctx, validSchema(name)); err != nil {
				t.Fatalf("CreateTable(%q) unexpected error: %v", name, err)
			}
		}
		names := e.Tables()
		if len(names) != 3 {
			t.Fatalf("Tables() returned %d names, want 3", len(names))
		}
		if names[0] != "alpha" || names[1] != "middle" || names[2] != "zeta" {
			t.Fatalf("Tables() = %v, want [alpha middle zeta]", names)
		}
	})

	t.Run("Tables/noTables/returnsEmpty", func(t *testing.T) {
		t.Parallel()
		e := New()
		if len(e.Tables()) != 0 {
			t.Fatalf("Tables() on empty engine returned %d, want 0", len(e.Tables()))
		}
	})
}

type stubStore struct{}

func (s *stubStore) CreateTable(context.Context, TableSchema) error               { return nil }
func (s *stubStore) DropTable(context.Context, string) error                      { return nil }
func (s *stubStore) Insert(context.Context, string, []Row) error                  { return nil }
func (s *stubStore) Update(context.Context, string, string, map[string]any) error { return nil }
func (s *stubStore) Delete(context.Context, string, string) error                 { return nil }
func (s *stubStore) Query(context.Context, string, Expr) (RowIterator, error)     { return nil, nil }
func (s *stubStore) Schema(context.Context, string) (TableSchema, error)          { return TableSchema{}, nil }
func (s *stubStore) Get(context.Context, string, string) (Row, error)             { return nil, nil }
