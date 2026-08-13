package service

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/domain"
	"backend/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, id int64, req domain.UpdateProductRequest) (*domain.Product, error)
	List(ctx context.Context, page, limit int, q string) (*domain.ProductListResponse, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{productRepo: repo}
}

func (s *productService) Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	req.SKU = strings.TrimSpace(req.SKU)
	req.Name = strings.TrimSpace(req.Name)

	if req.SKU == "" || req.Name == "" {
		return nil, domain.NewBadRequestError("invalid request", "sku and name are required")
	}
	if req.Price <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "price must be greater than 0")
	}
	if req.Stock < 0 {
		return nil, domain.NewBadRequestError("invalid request", "stock cannot be negative")
	}

	product := &domain.Product{
		SKU:   req.SKU,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "id must be positive")
	}
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, domain.NewNotFoundError("product not found", fmt.Sprintf("product with id %d does not exist", id))
	}
	return product, nil
}

func (s *productService) Update(ctx context.Context, id int64, req domain.UpdateProductRequest) (*domain.Product, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "id must be positive")
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, domain.NewBadRequestError("invalid request", "name is required")
	}
	if req.Price <= 0 {
		return nil, domain.NewBadRequestError("invalid request", "price must be greater than 0")
	}
	if req.Stock < 0 {
		return nil, domain.NewBadRequestError("invalid request", "stock cannot be negative")
	}

	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.NewNotFoundError("product not found", fmt.Sprintf("product with id %d does not exist", id))
	}

	existing.Name = req.Name
	existing.Price = req.Price
	existing.Stock = req.Stock

	if err := s.productRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *productService) List(ctx context.Context, page, limit int, q string) (*domain.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, totalData, err := s.productRepo.List(ctx, page, limit, q)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalData + int64(limit) - 1) / int64(limit))
	if totalPages == 0 && totalData == 0 {
		totalPages = 0
	}

	return &domain.ProductListResponse{
		Data: products,
		Pagination: domain.PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalData:  totalData,
			TotalPages: totalPages,
		},
	}, nil
}
