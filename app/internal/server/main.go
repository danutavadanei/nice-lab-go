package server

import (
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/gorilla/handlers"
	"net/http"
)

func StartHttpServer(
	cfg config.HTTPServerConfig,
	router http.Handler,
	shutdown chan bool,
	opts ...handlers.CORSOption,
) *http.Server {
	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handlers.CORS(opts...)(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			close(shutdown)
		}
	}()

	return srv
}
