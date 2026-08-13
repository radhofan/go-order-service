package handler

import (
	"net/http"
	"strconv"

	"backend/internal/domain"
	"backend/internal/handler/helper"
	"backend/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, err)
		return
	}

	product, err := h.productService.Create(r.Context(), req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	q := r.URL.Query().Get("q")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	res, err := h.productService.List(r.Context(), page, limit, q)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, res)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.WriteError(w, domain.NewBadRequestError("invalid product id", "id must be an integer"))
		return
	}

	product, err := h.productService.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.WriteError(w, domain.NewBadRequestError("invalid product id", "id must be an integer"))
		return
	}

	var req domain.UpdateProductRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, err)
		return
	}

	product, err := h.productService.Update(r.Context(), id, req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, product)
}
