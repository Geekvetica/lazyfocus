//go:build darwin

package bridge

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestExecute_SuccessfulExecution tests that successful osascript execution returns stdout
func TestExecute_SuccessfulExecution(t *testing.T) {
	executor := NewOSAScriptExecutor()

	// Simple JavaScript that returns a value
	script := `(() => { return "hello"; })()`

	result, err := executor.Execute(script)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// osascript returns the value with a newline
	expected := "hello\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestExecute_NonZeroExitCode tests that non-zero exit codes return error with stderr
func TestExecute_NonZeroExitCode(t *testing.T) {
	executor := NewOSAScriptExecutor()

	// Invalid JavaScript that will cause osascript to fail
	script := `(() => { throw new Error("test error"); })()`

	_, err := executor.Execute(script)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Error message should contain information about the failure
	if !strings.Contains(err.Error(), "execution failed") {
		t.Errorf("error should mention execution failure, got: %v", err)
	}
}

// TestExecuteWithTimeout_Success tests successful execution with custom timeout
func TestExecuteWithTimeout_Success(t *testing.T) {
	executor := NewOSAScriptExecutor()

	script := `(() => { return "timeout test"; })()`
	timeout := 5 * time.Second

	result, err := executor.ExecuteWithTimeout(script, timeout)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	expected := "timeout test\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestExecuteWithTimeout_TimeoutOccurs tests that timeout is properly enforced
func TestExecuteWithTimeout_TimeoutOccurs(t *testing.T) {
	executor := NewOSAScriptExecutor()

	// Script that delays longer than timeout
	script := `(() => {
		const start = Date.now();
		while (Date.now() - start < 3000) {}
		return "should not complete";
	})()`

	timeout := 100 * time.Millisecond

	_, err := executor.ExecuteWithTimeout(script, timeout)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !errors.Is(err, ErrExecutionTimeout) {
		t.Errorf("expected ErrExecutionTimeout, got: %v", err)
	}
}

// TestExecute_UsesDefaultTimeout tests that Execute uses the default timeout
func TestExecute_UsesDefaultTimeout(t *testing.T) {
	// This is implicitly tested by other tests, but we verify the behavior
	executor := NewOSAScriptExecutorWithTimeout(1 * time.Millisecond)

	// Script that takes longer than 1ms
	script := `(() => {
		const start = Date.now();
		while (Date.now() - start < 100) {}
		return "slow";
	})()`

	_, err := executor.Execute(script)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !errors.Is(err, ErrExecutionTimeout) {
		t.Errorf("expected ErrExecutionTimeout, got: %v", err)
	}
}

// TestExecute_EmptyScript tests handling of empty script
func TestExecute_EmptyScript(t *testing.T) {
	executor := NewOSAScriptExecutor()

	result, err := executor.Execute("")

	// osascript with empty script returns empty output
	if err != nil {
		t.Errorf("expected no error for empty script, got: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
}

// TestExecute_ContextCancellation tests that context cancellation works
func TestExecute_ContextCancellation(t *testing.T) {
	executor := NewOSAScriptExecutor()

	// Long-running script
	script := `(() => {
		const start = Date.now();
		while (Date.now() - start < 5000) {}
		return "completed";
	})()`

	// Very short timeout to trigger cancellation
	timeout := 10 * time.Millisecond

	_, err := executor.ExecuteWithTimeout(script, timeout)

	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}

	// Should be timeout error
	if !errors.Is(err, ErrExecutionTimeout) {
		t.Errorf("expected ErrExecutionTimeout, got: %v", err)
	}
}
