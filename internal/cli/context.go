package cli

import (
	"context"
	"errors"
	"reflect"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
)

// contextKey is used to store values in context
type contextKey string

const serviceKey contextKey = "omnifocus-service"

// ErrServiceNotFound is returned when service is not found in context
var ErrServiceNotFound = errors.New("service not found in context")

// isTypedNil detects typed-nil values (non-nil interface wrapping nil pointer)
func isTypedNil(v any) bool {
	if v == nil {
		return false // plain nil, not typed-nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// ServiceFromContext extracts the OmniFocusService from the context.
// Returns ErrServiceNotFound if the context is nil, the service is not present,
// or if the service is a typed-nil (non-nil interface wrapping a nil pointer).
func ServiceFromContext(ctx context.Context) (service.OmniFocusService, error) {
	// Guard against nil context to prevent panic on ctx.Value() call.
	// This can happen if a command is executed without ExecuteContext
	// and without the root PersistentPreRunE setting up the context.
	if ctx == nil {
		return nil, ErrServiceNotFound
	}
	svc, ok := ctx.Value(serviceKey).(service.OmniFocusService)
	// Check for both plain nil and typed-nil implementations.
	// A typed-nil occurs when a nil pointer is passed as an interface value
	// (e.g., var s *DefaultOmniFocusService = nil passed as OmniFocusService).
	if !ok || svc == nil || isTypedNil(svc) {
		return nil, ErrServiceNotFound
	}
	return svc, nil
}

// ContextWithService returns a new context with the service attached.
// If ctx is nil, context.Background() is used as the parent context
// to prevent panics from context.WithValue().
func ContextWithService(ctx context.Context, svc service.OmniFocusService) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, serviceKey, svc)
}
