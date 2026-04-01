package eval

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type mockInvoker struct {
	fn func(ctx context.Context, udfName string, inputs map[string]any) (map[string]any, error)
}

func (m *mockInvoker) Invoke(ctx context.Context, udfName string, inputs map[string]any) (map[string]any, error) {
	return m.fn(ctx, udfName, inputs)
}

func TestEvalSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		udf     string
		inputs  map[string]any
		want    map[string]any
	}{
		{
			"Eval/simpleUDF/returnsResult",
			"transcribe",
			map[string]any{"audio": "file.mp3"},
			map[string]any{"text": "hello world"},
		},
		{
			"Eval/emptyInputs/succeeds",
			"noop",
			nil,
			map[string]any{"status": "ok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			invoker := &mockInvoker{fn: func(_ context.Context, _ string, _ map[string]any) (map[string]any, error) {
				return tt.want, nil
			}}
			ev := New(invoker, Config{MaxRetries: 0, Timeout: 5 * time.Second})
			got, err := ev.Eval(context.Background(), tt.udf, tt.inputs)
			if err != nil {
				t.Fatalf("Eval() unexpected error: %v", err)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Fatalf("Eval() result[%q] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func TestEvalRetryThenSuccess(t *testing.T) {
	t.Parallel()

	t.Run("Eval/failsThenSucceeds/retriesAndReturnsResult", func(t *testing.T) {
		t.Parallel()
		var calls atomic.Int32
		invoker := &mockInvoker{fn: func(_ context.Context, _ string, _ map[string]any) (map[string]any, error) {
			n := calls.Add(1)
			if n < 3 {
				return nil, errors.New("transient failure")
			}
			return map[string]any{"done": true}, nil
		}}

		ev := New(invoker, Config{
			MaxRetries: 3,
			RetryDelay: time.Millisecond,
			Timeout:    5 * time.Second,
		})

		got, err := ev.Eval(context.Background(), "flaky", nil)
		if err != nil {
			t.Fatalf("Eval() unexpected error: %v", err)
		}
		if got["done"] != true {
			t.Fatalf("Eval() result = %v, want done=true", got)
		}
		if calls.Load() != 3 {
			t.Fatalf("expected 3 invocations, got %d", calls.Load())
		}
	})
}

func TestEvalAllRetriesExhausted(t *testing.T) {
	t.Parallel()

	t.Run("Eval/allFail/returnsLastError", func(t *testing.T) {
		t.Parallel()
		invokeErr := errors.New("permanent failure")
		var calls atomic.Int32
		invoker := &mockInvoker{fn: func(_ context.Context, _ string, _ map[string]any) (map[string]any, error) {
			calls.Add(1)
			return nil, invokeErr
		}}

		ev := New(invoker, Config{
			MaxRetries: 2,
			RetryDelay: time.Millisecond,
			Timeout:    5 * time.Second,
		})

		_, err := ev.Eval(context.Background(), "broken", nil)
		if err == nil {
			t.Fatalf("Eval() expected error, got nil")
		}
		if !errors.Is(err, invokeErr) {
			t.Fatalf("Eval() error = %v, want wrapping %v", err, invokeErr)
		}
		if calls.Load() != 3 {
			t.Fatalf("expected 3 total attempts (1 + 2 retries), got %d", calls.Load())
		}
	})
}

func TestEvalTimeout(t *testing.T) {
	t.Parallel()

	t.Run("Eval/excedsTimeout/returnsContextError", func(t *testing.T) {
		t.Parallel()
		invoker := &mockInvoker{fn: func(ctx context.Context, _ string, _ map[string]any) (map[string]any, error) {
			select {
			case <-time.After(5 * time.Second):
				return map[string]any{"late": true}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}}

		ev := New(invoker, Config{
			MaxRetries: 0,
			Timeout:    50 * time.Millisecond,
		})

		_, err := ev.Eval(context.Background(), "slow", nil)
		if err == nil {
			t.Fatalf("Eval() expected timeout error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Eval() error = %v, want context.DeadlineExceeded", err)
		}
	})
}
