package handler

import (
	"net/http"

	"backend/internal/domain"
	"backend/internal/handler/helper"
	"backend/internal/service"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(s service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: s}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCustomerRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, err)
		return
	}

	customer, err := h.customerService.Create(r.Context(), req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, customer)
}
