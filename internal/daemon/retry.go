package daemon

import (
	"context"
	"time"
)

// retryPolicy defines how scan errors are retried within a single tick.
type retryPolicy struct {
	maxAttempts int
	baseDelay   time.Duration
	maxDelay    time.Duration
}

type retryResult struct {
	attempts int
	err      error
}

func newRetryPolicy(maxAttempts int, baseDelay, maxDelay time.Duration) retryPolicy {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if baseDelay <= 0 {
		baseDelay = 100 * time.Millisecond
	}
	if maxDelay <= 0 || maxDelay < baseDelay {
		maxDelay = baseDelay * 8
	}
	return retryPolicy{
		maxAttempts: maxAttempts,
		baseDelay:   baseDelay,
		maxDelay:    maxDelay,
	}
}

// run executes fn up to maxAttempts times, backing off between failures.
// It returns early if ctx is cancelled or fn succeeds.
func (r retryPolicy) run(ctx context.Context, fn func() error) retryResult {
	var err error
	delay := r.baseDelay

	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		if ctx.Err() != nil {
			return retryResult{attempts: attempt - 1, err: ctx.Err()}
		}

		err = fn()
		if err == nil {
			return retryResult{attempts: attempt, err: nil}
		}

		if attempt == r.maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return retryResult{attempts: attempt, err: ctx.Err()}
		case <-time.After(delay):
		}

		delay *= 2
		if delay > r.maxDelay {
			delay = r.maxDelay
		}
	}

	return retryResult{attempts: r.maxAttempts, err: err}
}
