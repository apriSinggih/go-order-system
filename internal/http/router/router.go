package router

import (
	"net/http"

	"go-order-system/internal/http/handler"
	"go-order-system/pkg/route"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(
	authH *handler.AuthHandler,
	productH *handler.ProductHandler,
	orderH *handler.OrderHandler,
	authMiddleware gin.HandlerFunc,
) []route.Route {
	return []route.Route{
		// Health Check
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: handler.HealthCheck,
		},
		// Auth Endpoints
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/auth/register",
			Handler: authH.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/auth/login",
			Handler: authH.Login,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/v1/auth/me",
			Handler:     authH.Me,
			Middlewares: []gin.HandlerFunc{authMiddleware},
		},
		// Product Endpoints
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/products",
			Handler: productH.GetAll,
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/products/:id",
			Handler: productH.GetByID,
		},
		// Order Endpoints (butuh login)
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/orders",
			Handler:     orderH.Create,
			Middlewares: []gin.HandlerFunc{authMiddleware},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/v1/orders/:id",
			Handler:     orderH.GetByID,
			Middlewares: []gin.HandlerFunc{authMiddleware},
		},
	}
}
