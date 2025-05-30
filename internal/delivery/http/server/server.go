package http_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FlyKarlik/effectiveMobile/config"
	http_router "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/router"
)

type HTTPServer struct {
	cfg        *config.Config
	router     *http_router.HTTPRouter
	httpserver *http.Server
}

func New(cfg *config.Config, router *http_router.HTTPRouter) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		router: router,
		httpserver: func() *http.Server {
			return &http.Server{
				Addr:           fmt.Sprintf("%s:%s", cfg.AppUsers.AppHost, cfg.AppUsers.AppPort),
				Handler:        router.InitRouter(),
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				IdleTimeout:    10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
		}(),
	}
}

func (h *HTTPServer) ListenAndServe() error {
	return h.httpserver.ListenAndServe()
}

func (h *HTTPServer) Shuttdown(ctx context.Context) error {
	return h.httpserver.Shutdown(ctx)
}
