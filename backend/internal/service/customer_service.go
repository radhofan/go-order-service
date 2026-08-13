package service

import (
	"context"
	"net/mail"
	"strings"

	"backend/internal/domain"
	"backend/internal/repository"
)

type CustomerService interface {
	Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error)
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
