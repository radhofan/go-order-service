package test

import (
	"context"
	"testing"

	"backend/internal/domain"
)

func TestResourceNotFoundReturns404(t *testing.T) {
	ctx := context.Background()
	orderSvc, productSvc, _ := setupTestDB(t)

	_, err := productSvc.GetByID(ctx, 99999)
	if err == nil {
		t.Fatalf("expected 404 error for non-existent product")
	}
	if appErr, ok := err.(*domain.AppError); !ok || appErr.Code != 404 {
		t.Errorf("expected 404 code, got %v", err)
	}

	_, err = orderSvc.GetByID(ctx, 99999)
	if err == nil {
		t.Fatalf("expected 404 error for non-existent order")
	}
	if appErr, ok := err.(*domain.AppError); !ok || appErr.Code != 404 {
		t.Errorf("expected 404 code, got %v", err)
	}
}
