package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
)

func TestServiceFromContext_WithValidService(t *testing.T) {
	// Arrange
	mockService := &service.MockOmniFocusService{}
	ctx := context.WithValue(context.Background(), serviceKey, mockService)

	// Act
	result, err := ServiceFromContext(ctx)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != mockService {
		t.Error("Expected to get the same service instance")
	}
}

func TestServiceFromContext_WithoutService(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result, err := ServiceFromContext(ctx)

	// Assert
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
	if result != nil {
		t.Error("Expected nil service when error occurs")
	}
}

func TestServiceFromContext_WithNilService(t *testing.T) {
	// Arrange
	ctx := context.WithValue(context.Background(), serviceKey, nil)

	// Act
	result, err := ServiceFromContext(ctx)

	// Assert
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
	if result != nil {
		t.Error("Expected nil service when error occurs")
	}
}

func TestContextWithService_ShouldCreateContextWithService(t *testing.T) {
	// Arrange
	mockService := &service.MockOmniFocusService{}
	ctx := context.Background()

	// Act
	newCtx := ContextWithService(ctx, mockService)

	// Assert
	retrievedService, err := ServiceFromContext(newCtx)
	if err != nil {
		t.Errorf("Expected no error when retrieving service, got %v", err)
	}
	if retrievedService != mockService {
		t.Error("Expected to retrieve the same service instance")
	}
}

func TestServiceFromContext_WithNilContext(t *testing.T) {
	// Act
	result, err := ServiceFromContext(nil) //nolint:staticcheck // intentionally testing nil context handling

	// Assert
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
	if result != nil {
		t.Error("Expected nil service when error occurs")
	}
}

func TestContextWithService_WithNilContext(t *testing.T) {
	// Arrange
	mockService := &service.MockOmniFocusService{}

	// Act
	newCtx := ContextWithService(nil, mockService) //nolint:staticcheck // intentionally testing nil context handling

	// Assert
	retrievedService, err := ServiceFromContext(newCtx)
	if err != nil {
		t.Errorf("Expected no error when retrieving service, got %v", err)
	}
	if retrievedService != mockService {
		t.Error("Expected to retrieve the same service instance")
	}
}

func TestServiceFromContext_WithTypedNilService(t *testing.T) {
	// Arrange
	var nilSvc *service.MockOmniFocusService = nil
	ctx := context.WithValue(context.Background(), serviceKey, nilSvc)

	// Act
	result, err := ServiceFromContext(ctx)

	// Assert
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
	if result != nil {
		t.Error("Expected nil service when error occurs")
	}
}
