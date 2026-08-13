package test

import (
	"path/filepath"
	"testing"

	"backend/internal/db"
	"backend/internal/repository"
	"backend/internal/service"
)

func setupTestDB(t *testing.T) (service.OrderService, service.ProductService, service.CustomerService) {
	t.Helper()

	schemaPath, err := filepath.Abs("../schema.sql")
	if err != nil {
		t.Fatalf("failed to resolve schema path: %v", err)
	}

	database, err := db.InitDB("sqlite", "file::memory:?mode=memory&cache=shared", schemaPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	t.Cleanup(func() {
		database.Close()
	})

	productRepo := repository.NewProductRepository(database)
	customerRepo := repository.NewCustomerRepository(database)
	orderRepo := repository.NewOrderRepository(database)

	productSvc := service.NewProductService(productRepo)
	customerSvc := service.NewCustomerService(customerRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo, customerRepo)

	return orderSvc, productSvc, customerSvc
}
