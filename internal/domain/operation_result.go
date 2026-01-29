package domain

// OperationResult represents the outcome of a write operation
type OperationResult struct {
	Success bool   // Whether the operation succeeded
	ID      string // ID of the affected task
	Message string // Human-readable message
}

// NewSuccessResult creates a successful result
func NewSuccessResult(id, message string) OperationResult {
	return OperationResult{
		Success: true,
		ID:      id,
		Message: message,
	}
}

// NewErrorResult creates a failed result
func NewErrorResult(message string) OperationResult {
	return OperationResult{
		Success: false,
		ID:      "",
		Message: message,
	}
}
