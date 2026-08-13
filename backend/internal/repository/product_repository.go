package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"backend/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	List(ctx context.Context, page, limit int, q string) ([]domain.Product, int64, error)
	DeductStockTx(ctx context.Context, tx *sql.Tx, id int64, qty int) error
	RestoreStockTx(ctx context.Context, tx *sql.Tx, id int64, qty int) error
}

type sqliteProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &sqliteProductRepository{db: db}
}

func (r *sqliteProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (sku, name, price, stock) VALUES (?, ?, ?, ?) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, p.SKU, p.Name, p.Price, p.Stock).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return domain.NewConflictError("duplicate sku", "product with sku '"+p.SKU+"' already exists")
		}
		return err
	}
	return nil
}

func (r *sqliteProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, sku, name, price, stock, created_at FROM products WHERE id = ?`
	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *sqliteProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `SELECT id, sku, name, price, stock, created_at FROM products WHERE sku = ?`
	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *sqliteProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.NewNotFoundError("product not found", "")
	}
	return nil
}

func (r *sqliteProductRepository) List(ctx context.Context, page, limit int, q string) ([]domain.Product, int64, error) {
	offset := (page - 1) * limit

	var countQuery string
	var listQuery string
	var args []interface{}

	if q != "" {
		countQuery = `SELECT COUNT(*) FROM products WHERE LOWER(name) LIKE ?`
		listQuery = `SELECT id, sku, name, price, stock, created_at FROM products WHERE LOWER(name) LIKE ? ORDER BY id DESC LIMIT ? OFFSET ?`
		pattern := "%" + strings.ToLower(q) + "%"
		args = append(args, pattern)
	} else {
		countQuery = `SELECT COUNT(*) FROM products`
		listQuery = `SELECT id, sku, name, price, stock, created_at FROM products ORDER BY id DESC LIMIT ? OFFSET ?`
	}

	var totalData int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	listArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	return products, totalData, nil
}

func (r *sqliteProductRepository) DeductStockTx(ctx context.Context, tx *sql.Tx, id int64, qty int) error {
	query := `UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?`
	res, err := tx.ExecContext(ctx, query, qty, id, qty)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("insufficient stock or product not found")
	}
	return nil
}

func (r *sqliteProductRepository) RestoreStockTx(ctx context.Context, tx *sql.Tx, id int64, qty int) error {
	query := `UPDATE products SET stock = stock + ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, qty, id)
	return err
}
