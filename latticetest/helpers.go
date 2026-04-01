package latticetest

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/i33ym/lattice"
)

type MockStore struct {
	CreateTableFunc func(ctx context.Context, schema lattice.TableSchema) error
	DropTableFunc   func(ctx context.Context, name string) error
	InsertFunc      func(ctx context.Context, table string, rows []lattice.Row) error
	UpdateFunc      func(ctx context.Context, table string, id string, cols map[string]any) error
	DeleteFunc      func(ctx context.Context, table string, id string) error
	QueryFunc       func(ctx context.Context, table string, filter lattice.Expr) (lattice.RowIterator, error)
	SchemaFunc      func(ctx context.Context, table string) (lattice.TableSchema, error)
	GetFunc         func(ctx context.Context, table string, id string) (lattice.Row, error)
}

func (m *MockStore) CreateTable(ctx context.Context, schema lattice.TableSchema) error {
	if m.CreateTableFunc != nil {
		return m.CreateTableFunc(ctx, schema)
	}
	return nil
}

func (m *MockStore) DropTable(ctx context.Context, name string) error {
	if m.DropTableFunc != nil {
		return m.DropTableFunc(ctx, name)
	}
	return nil
}

func (m *MockStore) Insert(ctx context.Context, table string, rows []lattice.Row) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, table, rows)
	}
	return nil
}

func (m *MockStore) Update(ctx context.Context, table string, id string, cols map[string]any) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, table, id, cols)
	}
	return nil
}

func (m *MockStore) Delete(ctx context.Context, table string, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, table, id)
	}
	return nil
}

func (m *MockStore) Query(ctx context.Context, table string, filter lattice.Expr) (lattice.RowIterator, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, table, filter)
	}
	return NewSliceIterator(nil), nil
}

func (m *MockStore) Schema(ctx context.Context, table string) (lattice.TableSchema, error) {
	if m.SchemaFunc != nil {
		return m.SchemaFunc(ctx, table)
	}
	return lattice.TableSchema{}, nil
}

func (m *MockStore) Get(ctx context.Context, table string, id string) (lattice.Row, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, table, id)
	}
	return nil, nil
}

type MockBlobStore struct {
	PutFunc    func(ctx context.Context, key string, r io.Reader, meta lattice.BlobMeta) error
	GetFunc    func(ctx context.Context, key string) (io.ReadCloser, lattice.BlobMeta, error)
	DeleteFunc func(ctx context.Context, key string) error
	ExistsFunc func(ctx context.Context, key string) (bool, error)
}

func (m *MockBlobStore) Put(ctx context.Context, key string, r io.Reader, meta lattice.BlobMeta) error {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, key, r, meta)
	}
	return nil
}

func (m *MockBlobStore) Get(ctx context.Context, key string) (io.ReadCloser, lattice.BlobMeta, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return io.NopCloser(strings.NewReader("")), lattice.BlobMeta{}, nil
}

func (m *MockBlobStore) Delete(ctx context.Context, key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key)
	}
	return nil
}

func (m *MockBlobStore) Exists(ctx context.Context, key string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, key)
	}
	return false, nil
}

type MockVectorStore struct {
	CreateCollectionFunc func(ctx context.Context, name string, dim int) error
	DropCollectionFunc   func(ctx context.Context, name string) error
	UpsertFunc           func(ctx context.Context, collection string, records []lattice.VectorRecord) error
	SearchFunc           func(ctx context.Context, collection string, vector []float32, k int) ([]lattice.VectorResult, error)
	DeleteFunc           func(ctx context.Context, collection string, ids []string) error
}

func (m *MockVectorStore) CreateCollection(ctx context.Context, name string, dim int) error {
	if m.CreateCollectionFunc != nil {
		return m.CreateCollectionFunc(ctx, name, dim)
	}
	return nil
}

func (m *MockVectorStore) DropCollection(ctx context.Context, name string) error {
	if m.DropCollectionFunc != nil {
		return m.DropCollectionFunc(ctx, name)
	}
	return nil
}

func (m *MockVectorStore) Upsert(ctx context.Context, collection string, records []lattice.VectorRecord) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, collection, records)
	}
	return nil
}

func (m *MockVectorStore) Search(ctx context.Context, collection string, vector []float32, k int) ([]lattice.VectorResult, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, collection, vector, k)
	}
	return nil, nil
}

func (m *MockVectorStore) Delete(ctx context.Context, collection string, ids []string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, collection, ids)
	}
	return nil
}

type SliceIterator struct {
	rows []lattice.Row
	pos  int
	err  error
}

func NewSliceIterator(rows []lattice.Row) *SliceIterator {
	return &SliceIterator{rows: rows, pos: -1}
}

func NewSliceIteratorWithError(rows []lattice.Row, err error) *SliceIterator {
	return &SliceIterator{rows: rows, pos: -1, err: err}
}

func (s *SliceIterator) Next() bool {
	if s.err != nil {
		return false
	}
	s.pos++
	return s.pos < len(s.rows)
}

func (s *SliceIterator) Row() lattice.Row {
	if s.pos < 0 || s.pos >= len(s.rows) {
		return nil
	}
	return s.rows[s.pos]
}

func (s *SliceIterator) Err() error {
	return s.err
}

func (s *SliceIterator) Close() error {
	return nil
}

func MustCreateEngine(t *testing.T, opts ...lattice.Option) *lattice.Engine {
	t.Helper()
	return lattice.New(opts...)
}

func SampleSchema() lattice.TableSchema {
	return lattice.TableSchema{
		Name: "podcasts",
		Columns: []lattice.ColumnSpec{
			{Name: "id", Type: lattice.StringType},
			{Name: "title", Type: lattice.StringType},
			{Name: "audio", Type: lattice.AudioType},
			{Name: "duration", Type: lattice.IntType},
			{Name: "published", Type: lattice.TimestampType},
		},
	}
}

func SampleRows(n int) []lattice.Row {
	rows := make([]lattice.Row, n)
	for i := range n {
		rows[i] = lattice.Row{
			"id":    fmt.Sprintf("row-%d", i),
			"audio": fmt.Sprintf("s3://bucket/episode-%d.mp3", i),
		}
	}
	return rows
}
