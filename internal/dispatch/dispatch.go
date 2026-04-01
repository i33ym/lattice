package dispatch

import "context"

// Task represents a unit of work to dispatch.
type Task struct {
	ID     string
	RowID  string
	Column string
	Func   func(ctx context.Context) error
}

// Result holds the outcome of a dispatched task.
type Result struct {
	TaskID string
	Err    error
}

// Dispatcher dispatches tasks for concurrent execution.
type Dispatcher interface {
	Dispatch(ctx context.Context, tasks []Task) <-chan Result
	Shutdown(ctx context.Context) error
}
