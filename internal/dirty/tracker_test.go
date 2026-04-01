package dirty

import (
	"sort"
	"sync"
	"testing"
)

func TestMarkAndIsDirty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		markRow   string
		markCol   string
		checkRow  string
		checkCol  string
		wantDirty bool
	}{
		{"Mark/sameCell/isDirty", "r1", "c1", "r1", "c1", true},
		{"Mark/differentRow/isClean", "r1", "c1", "r2", "c1", false},
		{"Mark/differentCol/isClean", "r1", "c1", "r1", "c2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := New()
			tr.Mark(tt.markRow, tt.markCol)
			if got := tr.IsDirty(tt.checkRow, tt.checkCol); got != tt.wantDirty {
				t.Fatalf("IsDirty(%q, %q) = %v, want %v", tt.checkRow, tt.checkCol, got, tt.wantDirty)
			}
		})
	}
}

func TestClear(t *testing.T) {
	t.Parallel()

	t.Run("Clear/markedCell/becomesClean", func(t *testing.T) {
		t.Parallel()
		tr := New()
		tr.Mark("r1", "c1")
		tr.Clear("r1", "c1")
		if tr.IsDirty("r1", "c1") {
			t.Fatalf("IsDirty after Clear = true, want false")
		}
	})

	t.Run("Clear/unmarkedCell/noEffect", func(t *testing.T) {
		t.Parallel()
		tr := New()
		tr.Clear("r1", "c1")
		if tr.IsDirty("r1", "c1") {
			t.Fatalf("IsDirty after Clear on unmarked = true, want false")
		}
	})
}

func TestDirtyCells(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		marks [][2]string
		want  int
	}{
		{"DirtyCells/noMarks/returnsEmpty", nil, 0},
		{"DirtyCells/twoMarks/returnsTwo", [][2]string{{"r1", "c1"}, {"r2", "c2"}}, 2},
		{"DirtyCells/duplicateMark/returnsOne", [][2]string{{"r1", "c1"}, {"r1", "c1"}}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := New()
			for _, m := range tt.marks {
				tr.Mark(m[0], m[1])
			}
			cells := tr.DirtyCells()
			if len(cells) != tt.want {
				t.Fatalf("DirtyCells() returned %d cells, want %d", len(cells), tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		marks [][2]string
		want  int
	}{
		{"Count/empty/returnsZero", nil, 0},
		{"Count/three/returnsThree", [][2]string{{"r1", "c1"}, {"r2", "c1"}, {"r1", "c2"}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := New()
			for _, m := range tt.marks {
				tr.Mark(m[0], m[1])
			}
			if got := tr.Count(); got != tt.want {
				t.Fatalf("Count() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReset(t *testing.T) {
	t.Parallel()

	t.Run("Reset/afterMarks/clearsAll", func(t *testing.T) {
		t.Parallel()
		tr := New()
		tr.Mark("r1", "c1")
		tr.Mark("r2", "c2")
		tr.Reset()
		if tr.Count() != 0 {
			t.Fatalf("Count() after Reset = %d, want 0", tr.Count())
		}
		if tr.IsDirty("r1", "c1") {
			t.Fatalf("IsDirty after Reset = true, want false")
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	t.Run("ConcurrentMarkAndClear/noRace/completesSuccessfully", func(t *testing.T) {
		t.Parallel()
		tr := New()
		var wg sync.WaitGroup

		for i := range 100 {
			wg.Add(2)
			row := "r"
			col := string(rune('a' + (i % 26)))
			go func() {
				defer wg.Done()
				tr.Mark(row, col)
			}()
			go func() {
				defer wg.Done()
				tr.Clear(row, col)
			}()
		}

		wg.Wait()

		cells := tr.DirtyCells()
		_ = cells

		count := tr.Count()
		if count < 0 {
			t.Fatalf("Count() = %d, should be non-negative", count)
		}
	})

	t.Run("ConcurrentMarkAndRead/noRace/completesSuccessfully", func(t *testing.T) {
		t.Parallel()
		tr := New()
		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(2)
			col := string(rune('a' + (i % 26)))
			go func() {
				defer wg.Done()
				tr.Mark("r1", col)
			}()
			go func() {
				defer wg.Done()
				_ = tr.IsDirty("r1", col)
			}()
		}

		wg.Wait()
	})
}

func TestDirtyColumns(t *testing.T) {
	t.Parallel()

	t.Run("DirtyColumns/multipleColumns/returnsAllForRow", func(t *testing.T) {
		t.Parallel()
		tr := New()
		tr.Mark("r1", "a")
		tr.Mark("r1", "b")
		tr.Mark("r2", "c")

		cols := tr.DirtyColumns("r1")
		sort.Strings(cols)
		if len(cols) != 2 || cols[0] != "a" || cols[1] != "b" {
			t.Fatalf("DirtyColumns(r1) = %v, want [a b]", cols)
		}
	})
}
