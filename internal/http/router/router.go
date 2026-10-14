package router

import (
	"net/http"

	"go-order-system/internal/http/handler"
	"go-order-system/pkg/route"
)

func PublicRoutes(productH *handler.ProductHandler, orderH *handler.OrderHandler) []route.Route {
	return []route.Route{
		// Health Check
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: handler.HealthCheck,
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
		// Order Endpoints
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/orders",
			Handler: orderH.Create,
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/orders/:id",
			Handler: orderH.GetByID,
		},
	}
}