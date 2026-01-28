// Package bridge provides the interface between LazyFocus and OmniFocus
// via Omni Automation (JavaScript for Automation). It handles script
// execution, timeout management, and response parsing.
package bridge

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Error types for executor operations
var (
	ErrOSAScriptNotFound   = errors.New("osascript not found")
	ErrExecutionTimeout    = errors.New("script execution timed out")
	ErrOmniFocusNotRunning = errors.New("OmniFocus is not running")
)

// Executor defines the interface for executing Omni Automation scripts
type Executor interface {
	Execute(script string) (string, error)
	ExecuteWithTimeout(script string, timeout time.Duration) (string, error)
}

// OSAScriptExecutor executes JavaScript via osascript command
type OSAScriptExecutor struct {
	timeout time.Duration
}

// NewOSAScriptExecutor creates a new executor with default 30s timeout
func NewOSAScriptExecutor() *OSAScriptExecutor {
	return &OSAScriptExecutor{
		timeout: 30 * time.Second,
	}
}

// NewOSAScriptExecutorWithTimeout creates a new executor with custom timeout
func NewOSAScriptExecutorWithTimeout(timeout time.Duration) *OSAScriptExecutor {
	return &OSAScriptExecutor{
		timeout: timeout,
	}
}

// Execute runs a JavaScript script via osascript using the default timeout
func (e *OSAScriptExecutor) Execute(script string) (string, error) {
	return e.ExecuteWithTimeout(script, e.timeout)
}

// ExecuteWithTimeout runs a JavaScript script via osascript with a custom timeout
func (e *OSAScriptExecutor) ExecuteWithTimeout(script string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript", "-e", script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Check if context was cancelled (timeout occurred)
	if ctx.Err() == context.DeadlineExceeded {
		return "", ErrExecutionTimeout
	}

	// Check if osascript command was not found
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// Non-zero exit code, include stderr in error
			return "", fmt.Errorf("osascript execution failed: %w: %s", err, stderr.String())
		}

		// Check if command not found
		if errors.Is(err, exec.ErrNotFound) {
			return "", ErrOSAScriptNotFound
		}

		// Other execution error
		return "", fmt.Errorf("failed to execute osascript: %w", err)
	}

	return stdout.String(), nil
}
