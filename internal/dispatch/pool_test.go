package dispatch

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		workers int
	}{
		{"NewPool/positiveWorkers/succeeds", 4},
		{"NewPool/zeroWorkers/usesDefault", 0},
		{"NewPool/negativeWorkers/usesDefault", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := NewPool(tt.workers)
			if p == nil {
				t.Fatalf("NewPool(%d) returned nil", tt.workers)
			}
		})
	}
}

func TestDispatchExecutesAllTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		taskCount int
	}{
		{"Dispatch/singleTask/completesOne", 1},
		{"Dispatch/multipleTasks/completesAll", 5},
		{"Dispatch/noTasks/closesChannelImmediately", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := NewPool(2)
			ctx := context.Background()

			var mu sync.Mutex
			var executed []string

			tasks := make([]Task, tt.taskCount)
			for i := range tt.taskCount {
				id := string(rune('a' + i))
				tasks[i] = Task{
					ID:    id,
					RowID: "r1",
					Func: func(_ context.Context) error {
						mu.Lock()
						executed = append(executed, id)
						mu.Unlock()
						return nil
					},
				}
			}

			results := p.Dispatch(ctx, tasks)
			var resultList []Result
			for r := range results {
				resultList = append(resultList, r)
			}

			if len(resultList) != tt.taskCount {
				t.Fatalf("got %d results, want %d", len(resultList), tt.taskCount)
			}

			for _, r := range resultList {
				if r.Err != nil {
					t.Fatalf("task %q returned unexpected error: %v", r.TaskID, r.Err)
				}
			}
		})
	}
}

func TestDispatchRespectsContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("Dispatch/cancelledContext/returnsContextError", func(t *testing.T) {
		t.Parallel()
		p := NewPool(1)

		p.sem <- struct{}{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tasks := []Task{
			{
				ID: "blocked",
				Func: func(_ context.Context) error {
					return nil
				},
			},
		}

		results := p.Dispatch(ctx, tasks)
		for r := range results {
			if r.Err == nil {
				t.Fatalf("expected context error, got nil")
			}
			if !errors.Is(r.Err, context.Canceled) {
				t.Fatalf("expected context.Canceled, got %v", r.Err)
			}
		}

		<-p.sem
	})
}

func TestDispatchWithFailingTasks(t *testing.T) {
	t.Parallel()

	t.Run("Dispatch/someFail/reportsErrors", func(t *testing.T) {
		t.Parallel()
		p := NewPool(4)
		ctx := context.Background()

		taskErr := errors.New("task failed")
		tasks := []Task{
			{ID: "ok", Func: func(_ context.Context) error { return nil }},
			{ID: "fail", Func: func(_ context.Context) error { return taskErr }},
			{ID: "ok2", Func: func(_ context.Context) error { return nil }},
		}

		results := p.Dispatch(ctx, tasks)
		var errs int
		for r := range results {
			if r.Err != nil {
				errs++
			}
		}
		if errs != 1 {
			t.Fatalf("expected 1 error, got %d", errs)
		}
	})
}

func TestDispatchConcurrent(t *testing.T) {
	t.Parallel()

	t.Run("Dispatch/concurrent/allComplete", func(t *testing.T) {
		t.Parallel()
		p := NewPool(4)
		ctx := context.Background()

		var mu sync.Mutex
		var ids []string

		tasks := make([]Task, 20)
		for i := range 20 {
			id := string(rune('A' + i))
			tasks[i] = Task{
				ID: id,
				Func: func(_ context.Context) error {
					time.Sleep(time.Millisecond)
					mu.Lock()
					ids = append(ids, id)
					mu.Unlock()
					return nil
				},
			}
		}

		results := p.Dispatch(ctx, tasks)
		count := 0
		for range results {
			count++
		}
		if count != 20 {
			t.Fatalf("got %d results, want 20", count)
		}

		mu.Lock()
		sort.Strings(ids)
		mu.Unlock()

		if len(ids) != 20 {
			t.Fatalf("executed %d tasks, want 20", len(ids))
		}
	})
}

func TestShutdown(t *testing.T) {
	t.Parallel()

	t.Run("Shutdown/pool/returnsNil", func(t *testing.T) {
		t.Parallel()
		p := NewPool(2)
		if err := p.Shutdown(context.Background()); err != nil {
			t.Fatalf("Shutdown() unexpected error: %v", err)
		}
	})
}
