package handler

import (
	"net/http"

	"go-order-system/internal/http/dto"
	"go-order-system/internal/http/middleware"
	"go-order-system/internal/service"
	"go-order-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Registrasi berhasil", gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Pengguna tidak terautentikasi")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Profil pengguna", user)
}

