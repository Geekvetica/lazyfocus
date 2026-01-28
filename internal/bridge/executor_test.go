package bridge

import (
	"testing"
	"time"
)

// TestNewOSAScriptExecutor_DefaultTimeout tests default timeout is set
func TestNewOSAScriptExecutor_DefaultTimeout(t *testing.T) {
	executor := NewOSAScriptExecutor()

	if executor.timeout != 30*time.Second {
		t.Errorf("expected default timeout of 30s, got: %v", executor.timeout)
	}
}

// TestNewOSAScriptExecutorWithTimeout_CustomTimeout tests custom timeout is set
func TestNewOSAScriptExecutorWithTimeout_CustomTimeout(t *testing.T) {
	customTimeout := 10 * time.Second
	executor := NewOSAScriptExecutorWithTimeout(customTimeout)

	if executor.timeout != customTimeout {
		t.Errorf("expected timeout of %v, got: %v", customTimeout, executor.timeout)
	}
}

// TestExecutor_Interface tests that OSAScriptExecutor implements Executor interface
func TestExecutor_Interface(t *testing.T) {
	var _ Executor = (*OSAScriptExecutor)(nil)
}
