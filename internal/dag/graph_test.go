package dag

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

func TestAddNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ids     []string
		wantErr bool
	}{
		{"AddNode/singleNode/succeeds", []string{"a"}, false},
		{"AddNode/multipleNodes/succeeds", []string{"a", "b", "c"}, false},
		{"AddNode/duplicate/returnsError", []string{"a", "a"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := New()
			var lastErr error
			for _, id := range tt.ids {
				lastErr = g.AddNode(id, "table", id)
			}
			if tt.wantErr && lastErr == nil {
				t.Fatalf("AddNode() expected error for duplicate, got nil")
			}
			if !tt.wantErr && lastErr != nil {
				t.Fatalf("AddNode() unexpected error: %v", lastErr)
			}
		})
	}
}

func TestHasNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		add    []string
		check  string
		expect bool
	}{
		{"HasNode/existing/returnsTrue", []string{"a", "b"}, "a", true},
		{"HasNode/missing/returnsFalse", []string{"a"}, "z", false},
		{"HasNode/emptyGraph/returnsFalse", nil, "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := New()
			for _, id := range tt.add {
				if err := g.AddNode(id, "t", id); err != nil {
					t.Fatalf("AddNode(%q) unexpected error: %v", id, err)
				}
			}
			if got := g.HasNode(tt.check); got != tt.expect {
				t.Fatalf("HasNode(%q) = %v, want %v", tt.check, got, tt.expect)
			}
		})
	}
}

func TestAddEdge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		nodes    []string
		from, to string
		wantErr  bool
	}{
		{"AddEdge/validNodes/succeeds", []string{"a", "b"}, "a", "b", false},
		{"AddEdge/missingFrom/returnsError", []string{"b"}, "a", "b", true},
		{"AddEdge/missingTo/returnsError", []string{"a"}, "a", "b", true},
		{"AddEdge/duplicateEdge/succeedsIdempotent", []string{"a", "b"}, "a", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := New()
			for _, id := range tt.nodes {
				if err := g.AddNode(id, "t", id); err != nil {
					t.Fatalf("AddNode(%q) unexpected error: %v", id, err)
				}
			}
			err := g.AddEdge(tt.from, tt.to)
			if tt.wantErr && err == nil {
				t.Fatalf("AddEdge(%q, %q) expected error, got nil", tt.from, tt.to)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("AddEdge(%q, %q) unexpected error: %v", tt.from, tt.to, err)
			}
		})
	}
}

func TestDependentsAndDependencies(t *testing.T) {
	t.Parallel()

	t.Run("Dependents/directDownstream/returnsTargets", func(t *testing.T) {
		t.Parallel()
		g := New()
		for _, id := range []string{"a", "b", "c"} {
			if err := g.AddNode(id, "t", id); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		if err := g.AddEdge("a", "b"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := g.AddEdge("a", "c"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		deps := g.Dependents("a")
		sort.Strings(deps)
		want := []string{"b", "c"}
		if !reflect.DeepEqual(deps, want) {
			t.Fatalf("Dependents(a) = %v, want %v", deps, want)
		}
	})

	t.Run("Dependencies/directUpstream/returnsSources", func(t *testing.T) {
		t.Parallel()
		g := New()
		for _, id := range []string{"a", "b", "c"} {
			if err := g.AddNode(id, "t", id); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		if err := g.AddEdge("a", "c"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := g.AddEdge("b", "c"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		deps := g.Dependencies("c")
		sort.Strings(deps)
		want := []string{"a", "b"}
		if !reflect.DeepEqual(deps, want) {
			t.Fatalf("Dependencies(c) = %v, want %v", deps, want)
		}
	})

	t.Run("Dependents/noDownstream/returnsEmpty", func(t *testing.T) {
		t.Parallel()
		g := New()
		if err := g.AddNode("a", "t", "a"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		deps := g.Dependents("a")
		if len(deps) != 0 {
			t.Fatalf("Dependents(a) = %v, want empty", deps)
		}
	})
}

func TestAllDependentsAndAllDependencies(t *testing.T) {
	t.Parallel()

	setup := func() *Graph {
		g := New()
		for _, id := range []string{"a", "b", "c", "d"} {
			if err := g.AddNode(id, "t", id); err != nil {
				panic(err)
			}
		}
		for _, e := range [][2]string{{"a", "b"}, {"b", "c"}, {"c", "d"}} {
			if err := g.AddEdge(e[0], e[1]); err != nil {
				panic(err)
			}
		}
		return g
	}

	t.Run("AllDependents/transitiveChain/returnsAll", func(t *testing.T) {
		t.Parallel()
		g := setup()
		all := g.AllDependents("a")
		sort.Strings(all)
		want := []string{"b", "c", "d"}
		if !reflect.DeepEqual(all, want) {
			t.Fatalf("AllDependents(a) = %v, want %v", all, want)
		}
	})

	t.Run("AllDependencies/transitiveChain/returnsAll", func(t *testing.T) {
		t.Parallel()
		g := setup()
		all := g.AllDependencies("d")
		sort.Strings(all)
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(all, want) {
			t.Fatalf("AllDependencies(d) = %v, want %v", all, want)
		}
	})

	t.Run("AllDependents/leafNode/returnsEmpty", func(t *testing.T) {
		t.Parallel()
		g := setup()
		all := g.AllDependents("d")
		if len(all) != 0 {
			t.Fatalf("AllDependents(d) = %v, want empty", all)
		}
	})
}

func TestSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nodes   []string
		edges   [][2]string
		wantErr bool
	}{
		{
			"Sort/linearChain/succeeds",
			[]string{"a", "b", "c"},
			[][2]string{{"a", "b"}, {"b", "c"}},
			false,
		},
		{
			"Sort/diamondDependency/succeeds",
			[]string{"a", "b", "c", "d"},
			[][2]string{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}},
			false,
		},
		{
			"Sort/singleNode/succeeds",
			[]string{"a"},
			nil,
			false,
		},
		{
			"Sort/emptyGraph/succeeds",
			nil,
			nil,
			false,
		},
		{
			"Sort/cycle/returnsCycleError",
			[]string{"a", "b", "c"},
			[][2]string{{"a", "b"}, {"b", "c"}, {"c", "a"}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := New()
			for _, id := range tt.nodes {
				if err := g.AddNode(id, "t", id); err != nil {
					t.Fatalf("AddNode(%q) unexpected error: %v", id, err)
				}
			}
			for _, e := range tt.edges {
				if err := g.AddEdge(e[0], e[1]); err != nil {
					t.Fatalf("AddEdge(%q, %q) unexpected error: %v", e[0], e[1], err)
				}
			}
			order, err := g.Sort()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Sort() expected error, got nil")
				}
				if !errors.Is(err, ErrCycleDetected) {
					t.Fatalf("Sort() error = %v, want errors.Is ErrCycleDetected", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Sort() unexpected error: %v", err)
			}
			if len(order) != len(tt.nodes) {
				t.Fatalf("Sort() returned %d nodes, want %d", len(order), len(tt.nodes))
			}
		})
	}
}

func TestSortDiamondOrder(t *testing.T) {
	t.Parallel()

	t.Run("Sort/diamond/aPrecedesBAndC_BAndCPrecedeD", func(t *testing.T) {
		t.Parallel()
		g := New()
		for _, id := range []string{"a", "b", "c", "d"} {
			if err := g.AddNode(id, "t", id); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		for _, e := range [][2]string{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}} {
			if err := g.AddEdge(e[0], e[1]); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		order, err := g.Sort()
		if err != nil {
			t.Fatalf("Sort() unexpected error: %v", err)
		}

		pos := make(map[string]int, len(order))
		for i, id := range order {
			pos[id] = i
		}
		if pos["a"] >= pos["b"] {
			t.Fatalf("a (pos %d) should come before b (pos %d)", pos["a"], pos["b"])
		}
		if pos["a"] >= pos["c"] {
			t.Fatalf("a (pos %d) should come before c (pos %d)", pos["a"], pos["c"])
		}
		if pos["b"] >= pos["d"] {
			t.Fatalf("b (pos %d) should come before d (pos %d)", pos["b"], pos["d"])
		}
		if pos["c"] >= pos["d"] {
			t.Fatalf("c (pos %d) should come before d (pos %d)", pos["c"], pos["d"])
		}
	})
}

func TestDetectCycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		nodes  []string
		edges  [][2]string
		expect bool
	}{
		{
			"DetectCycle/noCycle/returnsFalse",
			[]string{"a", "b", "c"},
			[][2]string{{"a", "b"}, {"b", "c"}},
			false,
		},
		{
			"DetectCycle/withCycle/returnsTrue",
			[]string{"a", "b", "c"},
			[][2]string{{"a", "b"}, {"b", "c"}, {"c", "a"}},
			true,
		},
		{
			"DetectCycle/emptyGraph/returnsFalse",
			nil,
			nil,
			false,
		},
		{
			"DetectCycle/singleNode/returnsFalse",
			[]string{"a"},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := New()
			for _, id := range tt.nodes {
				if err := g.AddNode(id, "t", id); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			for _, e := range tt.edges {
				if err := g.AddEdge(e[0], e[1]); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if got := g.DetectCycle(); got != tt.expect {
				t.Fatalf("DetectCycle() = %v, want %v", got, tt.expect)
			}
		})
	}
}
