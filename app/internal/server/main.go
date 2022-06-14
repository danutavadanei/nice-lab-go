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
	var handler http.Handler

	if len(opts) > 0 {
		handler = handlers.CORS(opts...)(router)
	} else {
		handler = router
	}

	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			close(shutdown)
		}
	}()

	return srv
}
