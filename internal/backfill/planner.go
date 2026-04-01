package backfill

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// WorkUnit represents a single backfill task for one cell.
type WorkUnit struct {
	RowID  string
	Column string
	UDF    string
	Inputs []string
}

// Plan represents a backfill plan for a column across a set of rows.
type Plan struct {
	Table  string
	Column string
	Units  []WorkUnit
}

// Planner creates backfill plans.
type Planner struct{}

// New creates a Planner.
func New() *Planner {
	return &Planner{}
}

// CreatePlan generates a backfill plan that creates one WorkUnit per row ID.
func (p *Planner) CreatePlan(table, column, udf string, inputs []string, rowIDs []string) Plan {
	units := make([]WorkUnit, len(rowIDs))
	inputsCopy := make([]string, len(inputs))
	copy(inputsCopy, inputs)

	for i, rowID := range rowIDs {
		units[i] = WorkUnit{
			RowID:  rowID,
			Column: column,
			UDF:    udf,
			Inputs: inputsCopy,
		}
	}

	return Plan{
		Table:  table,
		Column: column,
		Units:  units,
	}
}

// Progress tracks backfill execution progress.
type Progress struct {
	Total     int
	Completed int
	Failed    int
}

// Executor runs backfill plans.
type Executor struct {
	evalFn func(ctx context.Context, unit WorkUnit) error
}

// NewExecutor creates an Executor that calls evalFn for each work unit.
func NewExecutor(evalFn func(ctx context.Context, unit WorkUnit) error) *Executor {
	return &Executor{evalFn: evalFn}
}

// Execute runs every work unit in the plan concurrently. It returns progress
// and the first error encountered, but processes all units regardless of
// individual failures.
func (ex *Executor) Execute(ctx context.Context, plan Plan) (Progress, error) {
	var (
		completed atomic.Int64
		failed    atomic.Int64
		firstErr  error
		errOnce   sync.Once
		wg        sync.WaitGroup
	)

	for _, unit := range plan.Units {
		wg.Add(1)
		go func(u WorkUnit) {
			defer wg.Done()
			if err := ex.evalFn(ctx, u); err != nil {
				failed.Add(1)
				errOnce.Do(func() {
					firstErr = fmt.Errorf("executing backfill for row %q column %q: %w", u.RowID, u.Column, err)
				})
				return
			}
			completed.Add(1)
		}(unit)
	}

	wg.Wait()

	return Progress{
		Total:     len(plan.Units),
		Completed: int(completed.Load()),
		Failed:    int(failed.Load()),
	}, firstErr
}
