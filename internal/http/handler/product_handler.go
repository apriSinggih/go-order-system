package handler

import (
	"net/http"
	"strconv"

	"go-order-system/internal/service"
	"go-order-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil data produk", products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil detail produk", product)
}