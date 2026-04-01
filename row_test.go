package lattice

import (
	"errors"
	"reflect"
	"testing"
)

type testSliceIterator struct {
	rows []Row
	pos  int
	err  error
}

func newTestSliceIterator(rows []Row) *testSliceIterator {
	return &testSliceIterator{rows: rows, pos: -1}
}

func newTestSliceIteratorWithError(rows []Row, err error) *testSliceIterator {
	return &testSliceIterator{rows: rows, pos: -1, err: err}
}

func (s *testSliceIterator) Next() bool {
	if s.err != nil {
		return false
	}
	s.pos++
	return s.pos < len(s.rows)
}

func (s *testSliceIterator) Row() Row {
	if s.pos < 0 || s.pos >= len(s.rows) {
		return nil
	}
	return s.rows[s.pos]
}

func (s *testSliceIterator) Err() error  { return s.err }
func (s *testSliceIterator) Close() error { return nil }

func TestCollectRows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		iter    RowIterator
		want    []Row
		wantErr bool
	}{
		{
			"CollectRows/multipleRows/returnsAll",
			newTestSliceIterator([]Row{
				{"id": "1", "name": "alice"},
				{"id": "2", "name": "bob"},
			}),
			[]Row{
				{"id": "1", "name": "alice"},
				{"id": "2", "name": "bob"},
			},
			false,
		},
		{
			"CollectRows/emptyIterator/returnsNilSlice",
			newTestSliceIterator(nil),
			nil,
			false,
		},
		{
			"CollectRows/iteratorWithError/returnsError",
			newTestSliceIteratorWithError(nil, errors.New("read failure")),
			nil,
			true,
		},
		{
			"CollectRows/singleRow/returnsOne",
			newTestSliceIterator([]Row{{"id": "only"}}),
			[]Row{{"id": "only"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := CollectRows(tt.iter)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("CollectRows() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("CollectRows() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("CollectRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
