package bridge

import (
	"errors"
	"testing"
	"time"
)

// mockExecutor is a test executor that can be configured to succeed or fail
type mockExecutor struct {
	// executeFunc is called when Execute is invoked
	executeFunc func(script string) (string, error)
	// executeWithTimeoutFunc is called when ExecuteWithTimeout is invoked
	executeWithTimeoutFunc func(script string, timeout time.Duration) (string, error)
}

func (m *mockExecutor) Execute(script string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(script)
	}
	return "", errors.New("not implemented")
}

func (m *mockExecutor) ExecuteWithTimeout(script string, timeout time.Duration) (string, error) {
	if m.executeWithTimeoutFunc != nil {
		return m.executeWithTimeoutFunc(script, timeout)
	}
	return "", errors.New("not implemented")
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts = 3, got %d", config.MaxAttempts)
	}

	if config.InitialWait != 100*time.Millisecond {
		t.Errorf("Expected InitialWait = 100ms, got %v", config.InitialWait)
	}

	if config.MaxWait != 2*time.Second {
		t.Errorf("Expected MaxWait = 2s, got %v", config.MaxWait)
	}
}

func TestRetryableExecutor_SuccessOnFirstAttempt(t *testing.T) {
	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			return "success", nil
		},
	}

	config := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	result, err := retryExecutor.ExecuteWithTimeout("test script", 5*time.Second)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %q", result)
	}
}

func TestRetryableExecutor_SuccessAfterRetries(t *testing.T) {
	attemptCount := 0
	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			attemptCount++
			if attemptCount < 3 {
				return "", ErrExecutionTimeout
			}
			return "success on third attempt", nil
		},
	}

	config := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	start := time.Now()
	result, err := retryExecutor.ExecuteWithTimeout("test script", 5*time.Second)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success on third attempt" {
		t.Errorf("Expected result 'success on third attempt', got %q", result)
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}

	// Verify exponential backoff occurred (should wait ~10ms + 20ms = ~30ms)
	minExpectedWait := 20 * time.Millisecond // Allow some slack for test execution
	if elapsed < minExpectedWait {
		t.Errorf("Expected at least %v wait time for retries, got %v", minExpectedWait, elapsed)
	}
}

func TestRetryableExecutor_FailureAfterMaxAttempts(t *testing.T) {
	attemptCount := 0
	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			attemptCount++
			return "", ErrExecutionTimeout
		},
	}

	config := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	result, err := retryExecutor.ExecuteWithTimeout("test script", 5*time.Second)

	if err != ErrExecutionTimeout {
		t.Errorf("Expected ErrExecutionTimeout, got %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}
}

func TestRetryableExecutor_NoRetryOnNonTimeoutError(t *testing.T) {
	attemptCount := 0
	otherError := errors.New("some other error")

	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			attemptCount++
			return "", otherError
		},
	}

	config := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	result, err := retryExecutor.ExecuteWithTimeout("test script", 5*time.Second)

	if err != otherError {
		t.Errorf("Expected otherError, got %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

	// Should only attempt once since it's not a timeout error
	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt (no retry), got %d", attemptCount)
	}
}

func TestRetryableExecutor_ExponentialBackoff(t *testing.T) {
	attemptCount := 0

	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			attemptCount++
			return "", ErrExecutionTimeout
		},
	}

	config := RetryConfig{
		MaxAttempts: 4,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     50 * time.Millisecond,
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	// Single execution should make 4 attempts with exponential backoff
	start := time.Now()
	_, err := retryExecutor.ExecuteWithTimeout("test script", 1*time.Second)
	elapsed := time.Since(start)

	if err != ErrExecutionTimeout {
		t.Errorf("Expected ErrExecutionTimeout, got %v", err)
	}

	// Should make exactly 4 attempts
	if attemptCount != 4 {
		t.Errorf("Expected 4 attempts, got %d", attemptCount)
	}

	// Total wait should be at least: 10ms + 20ms + 40ms = 70ms
	minExpectedWait := 60 * time.Millisecond // Allow some slack
	if elapsed < minExpectedWait {
		t.Errorf("Expected at least %v for exponential backoff, got %v", minExpectedWait, elapsed)
	}
}

func TestRetryableExecutor_Execute(t *testing.T) {
	// Test that Execute method uses ExecuteWithTimeout with default timeout
	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			if timeout != 30*time.Second {
				t.Errorf("Expected default timeout of 30s, got %v", timeout)
			}
			return "success", nil
		},
	}

	config := DefaultRetryConfig()
	retryExecutor := NewRetryableExecutor(mock, config)

	result, err := retryExecutor.Execute("test script")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %q", result)
	}
}

func TestRetryableExecutor_MaxWaitCap(t *testing.T) {
	attemptCount := 0
	mock := &mockExecutor{
		executeWithTimeoutFunc: func(script string, timeout time.Duration) (string, error) {
			attemptCount++
			return "", ErrExecutionTimeout
		},
	}

	// Configure with MaxWait less than what exponential backoff would reach
	config := RetryConfig{
		MaxAttempts: 5,
		InitialWait: 50 * time.Millisecond,
		MaxWait:     80 * time.Millisecond, // Cap at 80ms (would otherwise reach 50, 100, 200, 400)
	}

	retryExecutor := NewRetryableExecutor(mock, config)

	start := time.Now()
	_, err := retryExecutor.ExecuteWithTimeout("test script", 5*time.Second)
	elapsed := time.Since(start)

	if err != ErrExecutionTimeout {
		t.Errorf("Expected ErrExecutionTimeout, got %v", err)
	}

	if attemptCount != 5 {
		t.Errorf("Expected 5 attempts, got %d", attemptCount)
	}

	// Total wait should be ~50 + 80 + 80 + 80 = ~290ms (with cap)
	// Without cap it would be: ~50 + 100 + 200 + 400 = ~750ms
	maxExpectedWait := 400 * time.Millisecond
	if elapsed > maxExpectedWait {
		t.Errorf("Expected elapsed time < %v (indicating MaxWait cap), got %v", maxExpectedWait, elapsed)
	}
}
