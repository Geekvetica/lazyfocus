package cli

import (
	"context"
	"errors"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
)

// contextKey is used to store values in context
type contextKey string

const serviceKey contextKey = "omnifocus-service"

// ErrServiceNotFound is returned when service is not found in context
var ErrServiceNotFound = errors.New("service not found in context")

// ServiceFromContext extracts the OmniFocusService from the context
func ServiceFromContext(ctx context.Context) (service.OmniFocusService, error) {
	svc, ok := ctx.Value(serviceKey).(service.OmniFocusService)
	if !ok || svc == nil {
		return nil, ErrServiceNotFound
	}
	return svc, nil
}

// ContextWithService returns a new context with the service attached
func ContextWithService(ctx context.Context, svc service.OmniFocusService) context.Context {
	return context.WithValue(ctx, serviceKey, svc)
}
