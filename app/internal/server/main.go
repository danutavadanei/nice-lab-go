package server

import (
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/gorilla/handlers"
	"net/http"
)

func StartHttpServer(cfg config.HTTPServerConfig, router http.Handler, shutdown chan bool) *http.Server {
	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"X-Session-Token"}),
		)(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			close(shutdown)
		}
	}()

	return srv
}
