package test

import (
	"context"
	"testing"

	"backend/internal/domain"
)

func TestInvalidStatusTransitionRejected(t *testing.T) {
	ctx := context.Background()
	orderSvc, productSvc, customerSvc := setupTestDB(t)

	cust, _ := customerSvc.Create(ctx, domain.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test4@example.com",
	})

	p1, _ := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-006",
		Name:  "Product 6",
		Price: 20.0,
		Stock: 10,
	})

	order, err := orderSvc.Create(ctx, domain.CreateOrderRequest{
		CustomerID: cust.ID,
		Items: []domain.OrderItemRequest{
			{ProductID: p1.ID, Qty: 1},
		},
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	_, err = orderSvc.UpdateStatus(ctx, order.ID, domain.UpdateOrderStatusRequest{
		Status: domain.StatusShipped,
	})

	if err == nil {
		t.Fatalf("expected error for invalid status transition, got nil")
	}

	appErr, ok := err.(*domain.AppError)
	if !ok || appErr.Code != 409 {
		t.Errorf("expected AppError code 409 Conflict, got: %v", err)
	}
}
