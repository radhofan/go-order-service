package test

import (
	"context"
	"testing"

	"backend/internal/domain"
)

func TestOrderCreationSuccessAndTotalCalculation(t *testing.T) {
	ctx := context.Background()
	orderSvc, productSvc, customerSvc := setupTestDB(t)

	cust, err := customerSvc.Create(ctx, domain.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	p1, err := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-001",
		Name:  "Product 1",
		Price: 100.0,
		Stock: 10,
	})
	if err != nil {
		t.Fatalf("failed to create product 1: %v", err)
	}

	p2, err := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-002",
		Name:  "Product 2",
		Price: 50.0,
		Stock: 5,
	})
	if err != nil {
		t.Fatalf("failed to create product 2: %v", err)
	}

	order, err := orderSvc.Create(ctx, domain.CreateOrderRequest{
		CustomerID: cust.ID,
		Items: []domain.OrderItemRequest{
			{ProductID: p1.ID, Qty: 2},
			{ProductID: p2.ID, Qty: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected order creation success, got error: %v", err)
	}

	expectedTotal := (2 * 100.0) + (1 * 50.0)
	if order.TotalAmount != expectedTotal {
		t.Errorf("expected total_amount %f, got %f", expectedTotal, order.TotalAmount)
	}

	p1Updated, _ := productSvc.GetByID(ctx, p1.ID)
	if p1Updated.Stock != 8 {
		t.Errorf("expected p1 stock to be 8, got %d", p1Updated.Stock)
	}

	p2Updated, _ := productSvc.GetByID(ctx, p2.ID)
	if p2Updated.Stock != 4 {
		t.Errorf("expected p2 stock to be 4, got %d", p2Updated.Stock)
	}
}

func TestOrderCreationFailsInsufficientStock(t *testing.T) {
	ctx := context.Background()
	orderSvc, productSvc, customerSvc := setupTestDB(t)

	cust, _ := customerSvc.Create(ctx, domain.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test2@example.com",
	})

	p1, _ := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-003",
		Name:  "Product 3",
		Price: 100.0,
		Stock: 2,
	})

	_, err := orderSvc.Create(ctx, domain.CreateOrderRequest{
		CustomerID: cust.ID,
		Items: []domain.OrderItemRequest{
			{ProductID: p1.ID, Qty: 5},
		},
	})

	if err == nil {
		t.Fatalf("expected error due to insufficient stock, got nil")
	}

	appErr, ok := err.(*domain.AppError)
	if !ok || appErr.Code != 400 {
		t.Errorf("expected AppError with status code 400, got: %v", err)
	}
}
