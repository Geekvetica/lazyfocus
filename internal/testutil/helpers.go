// Package testutil provides shared test utilities
package testutil

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}
