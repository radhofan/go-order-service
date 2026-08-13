package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"
)

func main() {
	cfg := config.Load()

	database, err := db.InitDB(cfg.DBDriver, cfg.DBSource, "schema.sql")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	productRepo := repository.NewProductRepository(database)
	customerRepo := repository.NewCustomerRepository(database)
	orderRepo := repository.NewOrderRepository(database)

	productSvc := service.NewProductService(productRepo)
	customerSvc := service.NewCustomerService(customerRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo, customerRepo)

	productHdl := handler.NewProductHandler(productSvc)
	customerHdl := handler.NewCustomerHandler(customerSvc)
	orderHdl := handler.NewOrderHandler(orderSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", productHdl.Create)
	mux.HandleFunc("GET /products", productHdl.List)
	mux.HandleFunc("GET /products/{id}", productHdl.GetByID)
	mux.HandleFunc("PUT /products/{id}", productHdl.Update)

	mux.HandleFunc("POST /customers", customerHdl.Create)

	mux.HandleFunc("POST /orders", orderHdl.Create)
	mux.HandleFunc("GET /orders", orderHdl.List)
	mux.HandleFunc("GET /orders/{id}", orderHdl.GetByID)
	mux.HandleFunc("PATCH /orders/{id}/status", orderHdl.UpdateStatus)

	handler := middleware.Logger(mux)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on port %s...", cfg.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
