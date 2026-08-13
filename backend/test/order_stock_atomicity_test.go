package test

import (
	"context"
	"testing"

	"backend/internal/domain"
)

func TestStockUnchangedAfterFailedOrderCreation(t *testing.T) {
	ctx := context.Background()
	orderSvc, productSvc, customerSvc := setupTestDB(t)

	cust, _ := customerSvc.Create(ctx, domain.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test3@example.com",
	})

	p1, _ := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-004",
		Name:  "Product 4",
		Price: 100.0,
		Stock: 10,
	})

	p2, _ := productSvc.Create(ctx, domain.CreateProductRequest{
		SKU:   "SKU-005",
		Name:  "Product 5",
		Price: 50.0,
		Stock: 1,
	})

	_, err := orderSvc.Create(ctx, domain.CreateOrderRequest{
		CustomerID: cust.ID,
		Items: []domain.OrderItemRequest{
			{ProductID: p1.ID, Qty: 2},
			{ProductID: p2.ID, Qty: 10},
		},
	})

	if err == nil {
		t.Fatalf("expected order creation to fail")
	}

	p1Check, _ := productSvc.GetByID(ctx, p1.ID)
	if p1Check.Stock != 10 {
		t.Errorf("atomicity broken: expected p1 stock to remain 10, got %d", p1Check.Stock)
	}

	p2Check, _ := productSvc.GetByID(ctx, p2.ID)
	if p2Check.Stock != 1 {
		t.Errorf("atomicity broken: expected p2 stock to remain 1, got %d", p2Check.Stock)
	}
}
