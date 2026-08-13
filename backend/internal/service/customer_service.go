package service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"backend/internal/domain"
	"backend/internal/repository"
)

type CustomerService interface {
	Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error)
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{customerRepo: repo}
}

func (s *customerService) Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" {
		return nil, domain.NewBadRequestError("invalid request", "name is required")
	}
	if req.Email == "" {
		return nil, domain.NewBadRequestError("invalid request", "email is required")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, domain.NewBadRequestError("invalid request", "invalid email address format")
	}

	customer := &domain.Customer{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "id must be positive")
	}
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.NewNotFoundError("customer not found", fmt.Sprintf("customer with id %d does not exist", id))
	}
	return customer, nil
}
