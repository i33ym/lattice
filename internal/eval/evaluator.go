package eval

import (
	"context"
	"fmt"
	"time"
)

// UDFInvoker is the interface for calling a user-defined function.
type UDFInvoker interface {
	Invoke(ctx context.Context, udfName string, inputs map[string]any) (map[string]any, error)
}

// Config holds evaluator configuration.
type Config struct {
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// DefaultConfig returns a Config with sensible defaults: 3 retries, 1s delay,
// and 30s timeout.
func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		RetryDelay: time.Second,
		Timeout:    30 * time.Second,
	}
}

// Evaluator orchestrates UDF invocation with retry and timeout logic.
type Evaluator struct {
	invoker UDFInvoker
	config  Config
}

// New creates an Evaluator with the given invoker and configuration.
func New(invoker UDFInvoker, cfg Config) *Evaluator {
	return &Evaluator{
		invoker: invoker,
		config:  cfg,
	}
}

// Eval evaluates a single UDF call with the given inputs. It retries on
// failure using linear backoff (delay * attempt) and respects the configured
// timeout.
func (e *Evaluator) Eval(ctx context.Context, udfName string, inputs map[string]any) (map[string]any, error) {
	if e.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()
	}

	var lastErr error
	attempts := e.config.MaxRetries + 1
	if attempts < 1 {
		attempts = 1
	}

	for attempt := 1; attempt <= attempts; attempt++ {
		result, err := e.invoker.Invoke(ctx, udfName, inputs)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempt < attempts {
			delay := e.config.RetryDelay * time.Duration(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, fmt.Errorf("evaluating UDF %q (attempt %d): %w", udfName, attempt, ctx.Err())
			}
		}
	}

	return nil, fmt.Errorf("evaluating UDF %q after %d attempts: %w", udfName, attempts, lastErr)
}
