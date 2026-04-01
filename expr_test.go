package lattice

import (
	"testing"
)

func TestExprTypesImplementInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr Expr
	}{
		{"And/implementsExpr", And{}},
		{"Or/implementsExpr", Or{}},
		{"Eq/implementsExpr", Eq{}},
		{"GT/implementsExpr", GT{}},
		{"LT/implementsExpr", LT{}},
		{"GTE/implementsExpr", GTE{}},
		{"LTE/implementsExpr", LTE{}},
		{"JSONPath/implementsExpr", JSONPath{}},
		{"VectorSearch/implementsExpr", VectorSearch{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var _ Expr = tt.expr
		})
	}
}

func TestNewAnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		exprs     []Expr
		wantCount int
	}{
		{"NewAnd/noExprs/emptySlice", nil, 0},
		{"NewAnd/singleExpr/oneElement", []Expr{Eq{Column: "a", Value: 1}}, 1},
		{"NewAnd/multipleExprs/allPresent", []Expr{
			Eq{Column: "a", Value: 1},
			GT{Column: "b", Value: 2},
		}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			and := NewAnd(tt.exprs...)
			if len(and.Exprs) != tt.wantCount {
				t.Fatalf("NewAnd() returned %d exprs, want %d", len(and.Exprs), tt.wantCount)
			}
		})
	}
}

func TestNewOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		exprs     []Expr
		wantCount int
	}{
		{"NewOr/noExprs/emptySlice", nil, 0},
		{"NewOr/singleExpr/oneElement", []Expr{Eq{Column: "x", Value: "y"}}, 1},
		{"NewOr/multipleExprs/allPresent", []Expr{
			LT{Column: "a", Value: 10},
			GTE{Column: "b", Value: 5},
			Eq{Column: "c", Value: "z"},
		}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			or := NewOr(tt.exprs...)
			if len(or.Exprs) != tt.wantCount {
				t.Fatalf("NewOr() returned %d exprs, want %d", len(or.Exprs), tt.wantCount)
			}
		})
	}
}

func TestExprConstructors(t *testing.T) {
	t.Parallel()

	t.Run("NewEq/setsFields/correctly", func(t *testing.T) {
		t.Parallel()
		eq := NewEq("name", "alice")
		if eq.Column != "name" || eq.Value != "alice" {
			t.Fatalf("NewEq() = %+v, want Column=name Value=alice", eq)
		}
	})

	t.Run("NewGT/setsFields/correctly", func(t *testing.T) {
		t.Parallel()
		gt := NewGT("age", 30)
		if gt.Column != "age" || gt.Value != 30 {
			t.Fatalf("NewGT() = %+v, want Column=age Value=30", gt)
		}
	})

	t.Run("NewLT/setsFields/correctly", func(t *testing.T) {
		t.Parallel()
		lt := NewLT("score", 100)
		if lt.Column != "score" || lt.Value != 100 {
			t.Fatalf("NewLT() = %+v, want Column=score Value=100", lt)
		}
	})

	t.Run("NewJSONPath/setsFields/correctly", func(t *testing.T) {
		t.Parallel()
		jp := NewJSONPath("meta", "$.tags[0]", "go")
		if jp.Column != "meta" || jp.Path != "$.tags[0]" || jp.Value != "go" {
			t.Fatalf("NewJSONPath() = %+v", jp)
		}
	})

	t.Run("NewVectorSearch/setsFields/correctly", func(t *testing.T) {
		t.Parallel()
		vs := NewVectorSearch("embedding", []float32{1.0, 2.0, 3.0}, 5)
		if vs.Column != "embedding" || len(vs.Vector) != 3 || vs.K != 5 {
			t.Fatalf("NewVectorSearch() = %+v", vs)
		}
	})
}
