package query

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		table string
	}{
		{"NewBuilder/withTableName/setsTable", "users"},
		{"NewBuilder/emptyName/setsEmpty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := NewBuilder(tt.table)
			q := b.Build()
			if q.Table != tt.table {
				t.Fatalf("Build().Table = %q, want %q", q.Table, tt.table)
			}
		})
	}
}

func TestWhereAddsFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		column     string
		op         Op
		value      any
		wantColumn string
		wantOp     Op
	}{
		{"Where/eqFilter/addsCorrectly", "name", OpEq, "alice", "name", OpEq},
		{"Where/gtFilter/addsCorrectly", "age", OpGT, 18, "age", OpGT},
		{"Where/ltFilter/addsCorrectly", "score", OpLT, 100, "score", OpLT},
		{"Where/gteFilter/addsCorrectly", "rank", OpGTE, 1, "rank", OpGTE},
		{"Where/lteFilter/addsCorrectly", "price", OpLTE, 9.99, "price", OpLTE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := NewBuilder("t").Where(tt.column, tt.op, tt.value).Build()
			if len(q.Filters) != 1 {
				t.Fatalf("Build().Filters has %d entries, want 1", len(q.Filters))
			}
			f := q.Filters[0]
			if f.Column != tt.wantColumn {
				t.Fatalf("Filter.Column = %q, want %q", f.Column, tt.wantColumn)
			}
			if f.Op != tt.wantOp {
				t.Fatalf("Filter.Op = %d, want %d", f.Op, tt.wantOp)
			}
			if f.Value != tt.value {
				t.Fatalf("Filter.Value = %v, want %v", f.Value, tt.value)
			}
		})
	}
}

func TestLimitAndOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{"LimitOffset/both/setsCorrectly", 10, 20, 10, 20},
		{"LimitOffset/onlyLimit/offsetZero", 5, 0, 5, 0},
		{"LimitOffset/neither/bothZero", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := NewBuilder("t")
			if tt.limit > 0 {
				b = b.Limit(tt.limit)
			}
			if tt.offset > 0 {
				b = b.Offset(tt.offset)
			}
			q := b.Build()
			if q.Limit != tt.wantLimit {
				t.Fatalf("Query.Limit = %d, want %d", q.Limit, tt.wantLimit)
			}
			if q.Offset != tt.wantOffset {
				t.Fatalf("Query.Offset = %d, want %d", q.Offset, tt.wantOffset)
			}
		})
	}
}

func TestBuildReturnsCompleteQuery(t *testing.T) {
	t.Parallel()

	t.Run("Build/fullQuery/allFieldsSet", func(t *testing.T) {
		t.Parallel()
		q := NewBuilder("events").
			Where("type", OpEq, "click").
			Where("count", OpGT, 5).
			Limit(50).
			Offset(10).
			Build()

		if q.Table != "events" {
			t.Fatalf("Table = %q, want %q", q.Table, "events")
		}
		if len(q.Filters) != 2 {
			t.Fatalf("Filters count = %d, want 2", len(q.Filters))
		}
		if q.Limit != 50 {
			t.Fatalf("Limit = %d, want 50", q.Limit)
		}
		if q.Offset != 10 {
			t.Fatalf("Offset = %d, want 10", q.Offset)
		}
	})
}

func TestMethodChaining(t *testing.T) {
	t.Parallel()

	t.Run("Chaining/multipleWheres/accumulatesFilters", func(t *testing.T) {
		t.Parallel()
		q := NewBuilder("logs").
			Where("level", OpEq, "error").
			Where("timestamp", OpGTE, 1000).
			Where("source", OpEq, "api").
			Build()

		if len(q.Filters) != 3 {
			t.Fatalf("expected 3 filters, got %d", len(q.Filters))
		}
		if q.Filters[0].Column != "level" {
			t.Fatalf("first filter column = %q, want %q", q.Filters[0].Column, "level")
		}
		if q.Filters[1].Column != "timestamp" {
			t.Fatalf("second filter column = %q, want %q", q.Filters[1].Column, "timestamp")
		}
		if q.Filters[2].Column != "source" {
			t.Fatalf("third filter column = %q, want %q", q.Filters[2].Column, "source")
		}
	})
}

func TestBuildCopiesFilters(t *testing.T) {
	t.Parallel()

	t.Run("Build/mutateAfterBuild/originalUnaffected", func(t *testing.T) {
		t.Parallel()
		b := NewBuilder("t").Where("a", OpEq, 1)
		q1 := b.Build()
		b.Where("b", OpEq, 2)
		q2 := b.Build()

		if len(q1.Filters) != 1 {
			t.Fatalf("q1 should have 1 filter, got %d", len(q1.Filters))
		}
		if len(q2.Filters) != 2 {
			t.Fatalf("q2 should have 2 filters, got %d", len(q2.Filters))
		}
	})
}
