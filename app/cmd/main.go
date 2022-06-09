package main

import (
	"context"
	"encoding/json"
	"github.com/danutavadanei/nice-lab-go/internal/adapters/mysql"
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/danutavadanei/nice-lab-go/internal/server"
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

	db, err := mysql.NewConnection(cfg.MySQLConfig)
	userRep := mysql.NewUserRepository(db)
	labRep := mysql.NewLabRepository(db)
	sessionRep := mysql.NewSessionRepository(db, userRep, labRep)
	authTokenRep := mysql.NewAuthTokenRepository(db, userRep)

	if err != nil {
		panic(err)
	}

	m := mux.NewRouter()

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET").Name("health")

	m.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := userRep.ListUsers(r.Context())

		if err != nil {
			log.Printf("error listing users:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(users)

		if err != nil {
			log.Printf("error marshaling users:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("GET").Name("listUsers")

	m.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			log.Printf("error parsing form:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email, password := r.FormValue("email"), r.FormValue("password")

		if err := userRep.CheckUserPassword(r.Context(), email, password); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := userRep.GetUserByEmail(r.Context(), email)

		if err != nil {
			log.Printf("error parsing form:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := authTokenRep.NewTokenForUserId(r.Context(), user.ID)

		if err != nil {
			log.Printf("error generating auth token:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(struct {
			User  mysql.User `json:"user"`
			Token string     `json:"token"`
		}{
			User:  user,
			Token: token,
		})

		if err != nil {
			log.Printf("error parsing form:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bytes)
	}).Methods("POST").Name("login")

	m.HandleFunc("/labs", func(w http.ResponseWriter, r *http.Request) {
		labs, err := labRep.ListLabs(r.Context())

		if err != nil {
			log.Printf("error listing labs:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(labs)

		if err != nil {
			log.Printf("error marshaling labs:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("GET").Name("listLabs")

	m.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		sessions, err := sessionRep.ListSessions(r.Context())

		if err != nil {
			log.Printf("error listing sessions:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(sessions)

		if err != nil {
			log.Printf("error marshaling sessions:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("GET").Name("listSessions")

	srvShutdown := make(chan bool)
	srv := server.StartHttpServer(cfg.HTTPServerConfig, m, srvShutdown)

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
