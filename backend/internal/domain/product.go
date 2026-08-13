package domain

import "time"

type Product struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalData  int64 `json:"total_data"`
	TotalPages int   `json:"total_pages"`
}

type ProductListResponse struct {
	Data       []Product      `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}
