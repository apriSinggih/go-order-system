package handler

import (
	"net/http"
	"strconv"

	"go-order-system/internal/http/dto"
	"go-order-system/internal/service"
	"go-order-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Pesanan berhasil dibuat", order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pesanan tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data pesanan", order)
}