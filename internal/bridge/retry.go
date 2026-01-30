package bridge

import (
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
}

// DefaultRetryConfig returns sensible defaults for retry behavior.
// - MaxAttempts: 3 retries
// - InitialWait: 100ms (first retry wait)
// - MaxWait: 2s (maximum wait between retries)
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     2 * time.Second,
	}
}

// RetryableExecutor wraps an executor with retry logic.
// It automatically retries failed operations with exponential backoff,
// but only for timeout errors. Other errors fail immediately.
type RetryableExecutor struct {
	executor Executor
	config   RetryConfig
}

// NewRetryableExecutor creates a new retryable executor with the given configuration.
func NewRetryableExecutor(executor Executor, config RetryConfig) *RetryableExecutor {
	return &RetryableExecutor{
		executor: executor,
		config:   config,
	}
}

// Execute runs the script with retry logic using the default 30-second timeout.
func (r *RetryableExecutor) Execute(script string) (string, error) {
	return r.ExecuteWithTimeout(script, 30*time.Second)
}

// ExecuteWithTimeout runs the script with retry logic and a custom timeout.
// Only timeout errors (ErrExecutionTimeout) are retried.
// Other errors are returned immediately without retry.
// Implements exponential backoff with a configurable maximum wait time.
func (r *RetryableExecutor) ExecuteWithTimeout(script string, timeout time.Duration) (string, error) {
	var lastErr error
	wait := r.config.InitialWait

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		result, err := r.executor.ExecuteWithTimeout(script, timeout)
		if err == nil {
			return result, nil
		}

		// Only retry on timeout errors
		if err != ErrExecutionTimeout {
			return "", err
		}

		lastErr = err

		// Don't wait after last attempt
		if attempt < r.config.MaxAttempts {
			time.Sleep(wait)
			// Exponential backoff
			wait *= 2
			if wait > r.config.MaxWait {
				wait = r.config.MaxWait
			}
		}
	}

	return "", lastErr
}
