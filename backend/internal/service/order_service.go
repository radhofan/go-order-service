package service

import (
	"context"
	"fmt"

	"backend/internal/domain"
	"backend/internal/repository"
)

type OrderService interface {
	Create(ctx context.Context, req domain.CreateOrderRequest) (*domain.OrderResponse, error)
	GetByID(ctx context.Context, id int64) (*domain.OrderResponse, error)
	List(ctx context.Context, customerID int64, status domain.OrderStatus) ([]domain.OrderResponse, error)
	UpdateStatus(ctx context.Context, id int64, req domain.UpdateOrderStatusRequest) (*domain.OrderResponse, error)
}

type orderService struct {
	orderRepo    repository.OrderRepository
	productRepo  repository.ProductRepository
	customerRepo repository.CustomerRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository,
) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

type resolvedItem struct {
	product *domain.Product
	qty     int
}

func (s *orderService) Create(ctx context.Context, req domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	if req.CustomerID <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "customer_id is required")
	}

	if len(req.Items) == 0 {
		return nil, domain.NewBadRequestError("invalid request", "items list cannot be empty")
	}

	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.NewNotFoundError("customer not found", fmt.Sprintf("customer with id %d does not exist", req.CustomerID))
	}

	resolvedItems := make([]resolvedItem, 0, len(req.Items))
	var totalAmount float64

	for _, itemReq := range req.Items {
		if itemReq.ProductID <= 0 {
			return nil, domain.NewBadRequestError("invalid request", "product_id must be positive")
		}
		if itemReq.Qty <= 0 {
			return nil, domain.NewBadRequestError("invalid request", "qty must be greater than 0")
		}

		prod, err := s.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}
		if prod == nil {
			return nil, domain.NewNotFoundError("product not found", fmt.Sprintf("product with id %d does not exist", itemReq.ProductID))
		}

		if prod.Stock < itemReq.Qty {
			return nil, domain.NewBadRequestError(
				"insufficient stock",
				fmt.Sprintf("product '%s' remaining %d, requested %d", prod.SKU, prod.Stock, itemReq.Qty),
			)
		}

		totalAmount += float64(itemReq.Qty) * prod.Price
		resolvedItems = append(resolvedItems, resolvedItem{
			product: prod,
			qty:     itemReq.Qty,
		})
	}

	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return nil, domain.NewInternalError("transaction error", err.Error())
	}
	defer tx.Rollback()

	orderID, err := s.orderRepo.CreateTx(ctx, tx, req.CustomerID, totalAmount)
	if err != nil {
		return nil, domain.NewInternalError("failed to create order", err.Error())
	}

	for _, item := range resolvedItems {
		if err := s.productRepo.DeductStockTx(ctx, tx, item.product.ID, item.qty); err != nil {
			return nil, domain.NewBadRequestError(
				"insufficient stock",
				fmt.Sprintf("product '%s' remaining stock insufficient", item.product.SKU),
			)
		}

		if err := s.orderRepo.CreateOrderItemTx(ctx, tx, orderID, item.product.ID, item.qty, item.product.Price); err != nil {
			return nil, domain.NewInternalError("failed to create order item", err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit transaction", err.Error())
	}

	return s.orderRepo.GetByID(ctx, orderID)
}

func (s *orderService) GetByID(ctx context.Context, id int64) (*domain.OrderResponse, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "id must be positive")
	}

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found", fmt.Sprintf("order with id %d does not exist", id))
	}
	return order, nil
}

func (s *orderService) List(ctx context.Context, customerID int64, status domain.OrderStatus) ([]domain.OrderResponse, error) {
	return s.orderRepo.List(ctx, customerID, status)
}

func (s *orderService) UpdateStatus(ctx context.Context, id int64, req domain.UpdateOrderStatusRequest) (*domain.OrderResponse, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "id must be positive")
	}

	currentOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentOrder == nil {
		return nil, domain.NewNotFoundError("order not found", fmt.Sprintf("order with id %d does not exist", id))
	}

	if !domain.IsValidStatusTransition(currentOrder.Status, req.Status) {
		return nil, domain.NewConflictError(
			"invalid status transition",
			fmt.Sprintf("cannot transition order status from %s to %s", currentOrder.Status, req.Status),
		)
	}

	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return nil, domain.NewInternalError("transaction error", err.Error())
	}
	defer tx.Rollback()

	if req.Status == domain.StatusCancelled {
		for _, item := range currentOrder.Items {
			if err := s.productRepo.RestoreStockTx(ctx, tx, item.ProductID, item.Qty); err != nil {
				return nil, domain.NewInternalError("failed to restore stock", err.Error())
			}
		}
	}

	if err := s.orderRepo.UpdateStatusTx(ctx, tx, id, req.Status); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit status update", err.Error())
	}

	return s.orderRepo.GetByID(ctx, id)
}
