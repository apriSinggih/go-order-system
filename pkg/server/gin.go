package server

import (
	"context"
	"log"
	"net/http"

	"go-order-system/config"
	"go-order-system/pkg/route"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	Logger     *log.Logger
}

func NewServer(cfg *config.Config, routes []route.Route) *Server {
	router := gin.Default()

	// Registrasi seluruh route yang dikirim dari Builder
	for _, r := range routes {
		router.Handle(r.Method, r.Path, r.Handler)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.AppPort,
			Handler: router,
		},
		Logger: log.Default(),
	}
}

func (s *Server) Start(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}