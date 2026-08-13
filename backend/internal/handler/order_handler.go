package handler

import (
	"net/http"
	"strconv"

	"backend/internal/domain"
	"backend/internal/handler/helper"
	"backend/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: s}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, err)
		return
	}

	order, err := h.orderService.Create(r.Context(), req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.WriteError(w, domain.NewBadRequestError("invalid order id", "id must be an integer"))
		return
	}

	order, err := h.orderService.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	statusStr := r.URL.Query().Get("status")

	var customerID int64
	if customerIDStr != "" {
		if cid, err := strconv.ParseInt(customerIDStr, 10, 64); err == nil {
			customerID = cid
		}
	}

	status := domain.OrderStatus(statusStr)

	orders, err := h.orderService.List(r.Context(), customerID, status)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.WriteError(w, domain.NewBadRequestError("invalid order id", "id must be an integer"))
		return
	}

	var req domain.UpdateOrderStatusRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, err)
		return
	}

	order, err := h.orderService.UpdateStatus(r.Context(), id, req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, order)
}
