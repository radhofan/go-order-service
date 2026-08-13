package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"backend/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *domain.Customer) error
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
}

type sqliteCustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &sqliteCustomerRepository{db: db}
}

func (r *sqliteCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (name, email) VALUES (?, ?) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, c.Name, c.Email).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return domain.NewConflictError("duplicate customer email", "customer with email '"+c.Email+"' already exists")
		}
		return err
	}
	return nil
}

func (r *sqliteCustomerRepository) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	query := `SELECT id, name, email, created_at FROM customers WHERE id = ?`
	c := &domain.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}
