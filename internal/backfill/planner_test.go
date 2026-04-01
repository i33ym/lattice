package backfill

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestCreatePlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		table     string
		column    string
		udf       string
		inputs    []string
		rowIDs    []string
		wantUnits int
	}{
		{
			"CreatePlan/multipleRows/generatesOnePerRow",
			"podcasts", "transcript", "transcribe",
			[]string{"audio"},
			[]string{"r1", "r2", "r3"},
			3,
		},
		{
			"CreatePlan/noRows/generatesEmptyPlan",
			"podcasts", "transcript", "transcribe",
			[]string{"audio"},
			nil,
			0,
		},
		{
			"CreatePlan/singleRow/generatesOneUnit",
			"events", "summary", "summarize",
			[]string{"payload", "title"},
			[]string{"r1"},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := New()
			plan := p.CreatePlan(tt.table, tt.column, tt.udf, tt.inputs, tt.rowIDs)

			if plan.Table != tt.table {
				t.Fatalf("Plan.Table = %q, want %q", plan.Table, tt.table)
			}
			if plan.Column != tt.column {
				t.Fatalf("Plan.Column = %q, want %q", plan.Column, tt.column)
			}
			if len(plan.Units) != tt.wantUnits {
				t.Fatalf("len(Plan.Units) = %d, want %d", len(plan.Units), tt.wantUnits)
			}

			for i, unit := range plan.Units {
				if unit.RowID != tt.rowIDs[i] {
					t.Fatalf("Unit[%d].RowID = %q, want %q", i, unit.RowID, tt.rowIDs[i])
				}
				if unit.Column != tt.column {
					t.Fatalf("Unit[%d].Column = %q, want %q", i, unit.Column, tt.column)
				}
				if unit.UDF != tt.udf {
					t.Fatalf("Unit[%d].UDF = %q, want %q", i, unit.UDF, tt.udf)
				}
			}
		})
	}
}

func TestCreatePlanInputsCopied(t *testing.T) {
	t.Parallel()

	t.Run("CreatePlan/inputsMutated/planUnaffected", func(t *testing.T) {
		t.Parallel()
		p := New()
		inputs := []string{"audio", "title"}
		plan := p.CreatePlan("t", "c", "fn", inputs, []string{"r1"})

		inputs[0] = "MUTATED"

		if plan.Units[0].Inputs[0] == "MUTATED" {
			t.Fatalf("CreatePlan did not copy inputs slice")
		}
	})
}

func TestExecuteAllSuccess(t *testing.T) {
	t.Parallel()

	t.Run("Execute/allSucceed/returnsFullProgress", func(t *testing.T) {
		t.Parallel()
		var count atomic.Int32
		exec := NewExecutor(func(_ context.Context, _ WorkUnit) error {
			count.Add(1)
			return nil
		})

		plan := Plan{
			Table:  "t",
			Column: "c",
			Units: []WorkUnit{
				{RowID: "r1", Column: "c", UDF: "fn"},
				{RowID: "r2", Column: "c", UDF: "fn"},
				{RowID: "r3", Column: "c", UDF: "fn"},
			},
		}

		progress, err := exec.Execute(context.Background(), plan)
		if err != nil {
			t.Fatalf("Execute() unexpected error: %v", err)
		}
		if progress.Total != 3 {
			t.Fatalf("Progress.Total = %d, want 3", progress.Total)
		}
		if progress.Completed != 3 {
			t.Fatalf("Progress.Completed = %d, want 3", progress.Completed)
		}
		if progress.Failed != 0 {
			t.Fatalf("Progress.Failed = %d, want 0", progress.Failed)
		}
	})
}

func TestExecuteWithFailures(t *testing.T) {
	t.Parallel()

	t.Run("Execute/someUnitsFailure/reportsPartialProgress", func(t *testing.T) {
		t.Parallel()
		unitErr := errors.New("compute error")
		exec := NewExecutor(func(_ context.Context, u WorkUnit) error {
			if u.RowID == "r2" {
				return unitErr
			}
			return nil
		})

		plan := Plan{
			Table:  "t",
			Column: "c",
			Units: []WorkUnit{
				{RowID: "r1", Column: "c", UDF: "fn"},
				{RowID: "r2", Column: "c", UDF: "fn"},
				{RowID: "r3", Column: "c", UDF: "fn"},
			},
		}

		progress, err := exec.Execute(context.Background(), plan)
		if err == nil {
			t.Fatalf("Execute() expected error, got nil")
		}
		if progress.Total != 3 {
			t.Fatalf("Progress.Total = %d, want 3", progress.Total)
		}
		if progress.Completed != 2 {
			t.Fatalf("Progress.Completed = %d, want 2", progress.Completed)
		}
		if progress.Failed != 1 {
			t.Fatalf("Progress.Failed = %d, want 1", progress.Failed)
		}
	})
}

func TestExecuteEmptyPlan(t *testing.T) {
	t.Parallel()

	t.Run("Execute/emptyPlan/returnsZeroProgress", func(t *testing.T) {
		t.Parallel()
		exec := NewExecutor(func(_ context.Context, _ WorkUnit) error {
			t.Fatalf("should not be called on empty plan")
			return nil
		})

		progress, err := exec.Execute(context.Background(), Plan{})
		if err != nil {
			t.Fatalf("Execute() unexpected error: %v", err)
		}
		if progress.Total != 0 || progress.Completed != 0 || progress.Failed != 0 {
			t.Fatalf("Progress = %+v, want all zeros", progress)
		}
	})
}
