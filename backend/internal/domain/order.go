package domain

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Qty       int   `json:"qty"`
}

type CreateOrderRequest struct {
	CustomerID int64              `json:"customer_id"`
	Items      []OrderItemRequest `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status"`
}

type OrderItemResponse struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"order_id"`
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name,omitempty"`
	Qty          int     `json:"qty"`
	PriceAtOrder float64 `json:"price_at_order"`
	Subtotal     float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID          int64               `json:"id"`
	CustomerID  int64               `json:"customer_id"`
	Status      OrderStatus         `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Items       []OrderItemResponse `json:"items,omitempty"`
}

func IsValidStatusTransition(current, next OrderStatus) bool {
	switch current {
	case StatusPending:
		return next == StatusPaid || next == StatusCancelled
	case StatusPaid:
		return next == StatusShipped || next == StatusCancelled
	case StatusShipped:
		return next == StatusCompleted
	default:
		return false
	}
}
