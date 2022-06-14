package main

import (
	"context"
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/danutavadanei/nice-lab-go/internal/proxy"
	"github.com/danutavadanei/nice-lab-go/internal/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	cfg := config.NewAppConfig(v)

	authProxy, err := proxy.NewProxy(cfg.GatewayConfig.AuthUrl)
	if err != nil {
		panic("Could not initialize auth proxy")
	}

	pipelineProxy, err := proxy.NewProxy(cfg.GatewayConfig.PipelineUrl)
	if err != nil {
		panic("Could not initialize pipeline proxy")
	}

	m := mux.NewRouter()
	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET").Name("health")

	m.HandleFunc("/v1/auth/{rest:.*}", proxy.ProxyRequestHandler("/v1/auth", authProxy))
	m.HandleFunc("/v1/pipeline/{rest:.*}", proxy.ProxyRequestHandler("/v1/pipeline", pipelineProxy))

	srvShutdown := make(chan bool)
	srv := server.StartHttpServer(
		cfg.HTTPServerConfig,
		m,
		srvShutdown,
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Session-Token"}),
	)

	<-sigChannel
	go shutdown(srv)
	<-srvShutdown
}

func shutdown(server *http.Server) {
	ctxShutDown, _ := context.WithTimeout(context.Background(), 30)
	err := server.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("error shutting down server (%s): %v", server.Addr, err)
		err = server.Close()
		if err != nil {
			log.Printf("error closing server (%s): %v", server.Addr, err)
		}
	}
}
