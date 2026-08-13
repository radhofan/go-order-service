package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTx(ctx context.Context, tx *sql.Tx, customerID int64, totalAmount float64) (int64, error)
	CreateOrderItemTx(ctx context.Context, tx *sql.Tx, orderID, productID int64, qty int, priceAtOrder float64) error
	GetByID(ctx context.Context, id int64) (*domain.OrderResponse, error)
	List(ctx context.Context, customerID int64, status domain.OrderStatus) ([]domain.OrderResponse, error)
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, orderID int64, status domain.OrderStatus) error
	GetOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItemResponse, error)
}

type sqliteOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &sqliteOrderRepository{db: db}
}

func (r *sqliteOrderRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *sqliteOrderRepository) CreateTx(ctx context.Context, tx *sql.Tx, customerID int64, totalAmount float64) (int64, error) {
	query := `INSERT INTO orders (customer_id, status, total_amount) VALUES (?, 'PENDING', ?) RETURNING id`
	var orderID int64
	err := tx.QueryRowContext(ctx, query, customerID, totalAmount).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *sqliteOrderRepository) CreateOrderItemTx(ctx context.Context, tx *sql.Tx, orderID, productID int64, qty int, priceAtOrder float64) error {
	query := `INSERT INTO order_items (order_id, product_id, qty, price_at_order) VALUES (?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, orderID, productID, qty, priceAtOrder)
	return err
}

func (r *sqliteOrderRepository) GetByID(ctx context.Context, id int64) (*domain.OrderResponse, error) {
	query := `SELECT id, customer_id, status, total_amount, created_at FROM orders WHERE id = ?`
	o := &domain.OrderResponse{}
	var statusStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.CustomerID, &statusStr, &o.TotalAmount, &o.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	o.Status = domain.OrderStatus(statusStr)

	items, err := r.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return o, nil
}

func (r *sqliteOrderRepository) List(ctx context.Context, customerID int64, status domain.OrderStatus) ([]domain.OrderResponse, error) {
	query := `SELECT id, customer_id, status, total_amount, created_at FROM orders WHERE 1=1`
	var args []interface{}

	if customerID > 0 {
		query += ` AND customer_id = ?`
		args = append(args, customerID)
	}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, string(status))
	}

	query += ` ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []domain.OrderResponse{}
	for rows.Next() {
		var o domain.OrderResponse
		var statusStr string
		if err := rows.Scan(&o.ID, &o.CustomerID, &statusStr, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.Status = domain.OrderStatus(statusStr)

		items, err := r.GetOrderItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *sqliteOrderRepository) UpdateStatusTx(ctx context.Context, tx *sql.Tx, orderID int64, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`
	res, err := tx.ExecContext(ctx, query, string(status), orderID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.NewNotFoundError("order not found", "")
	}
	return nil
}

func (r *sqliteOrderRepository) GetOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItemResponse, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, p.name, oi.qty, oi.price_at_order 
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.OrderItemResponse{}
	for rows.Next() {
		var item domain.OrderItemResponse
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Qty, &item.PriceAtOrder); err != nil {
			return nil, err
		}
		item.Subtotal = float64(item.Qty) * item.PriceAtOrder
		items = append(items, item)
	}

	return items, nil
}
