package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/danutavadanei/nice-lab-go/internal/adapters/mysql"
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/danutavadanei/nice-lab-go/internal/server"
	"github.com/danutavadanei/nice-lab-go/internal/server/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	v := viper.New()
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	// v.AutomaticEnv()

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	cfg := config.NewAppConfig(v)

	db, err := mysql.NewConnection(cfg.MySQLConfig)
	if err != nil {
		panic(err)
	}

	userRep := mysql.NewUserRepository(db)
	labRep := mysql.NewLabRepository(db)
	sessionRep := mysql.NewSessionRepository(db, userRep, labRep)
	authTokenRep := mysql.NewAuthTokenRepository(db, userRep)
	tokenUsers, err := authTokenRep.ListTokenUsers(context.Background())
	if err != nil {
		panic(err)
	}

	authMiddleware := middleware.NewAuthenticationMiddleware(tokenUsers, authTokenRep)

	ssmClient := ssm.NewFromConfig(*cfg.AWSConfig)

	m := mux.NewRouter()
	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET").Name("health")
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

	a := m.PathPrefix("/v1").Subrouter()
	a.Use(authMiddleware.Middleware)
	a.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		user := (r.Context().Value("user")).(mysql.User)

		if user.Type != mysql.Professor {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
	a.HandleFunc("/labs", func(w http.ResponseWriter, r *http.Request) {
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
	a.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		user := (r.Context().Value("user")).(mysql.User)

		if user.Type != mysql.Professor {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
	a.HandleFunc("/labs/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)

		if err != nil {
			log.Printf("error parsing request:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lab, err := labRep.GetLabById(r.Context(), id)

		if err != nil {
			log.Printf("error fetching lab:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := sessionRep.CreateSession(r.Context(), r.Context().Value("user").(mysql.User), lab)

		if err != nil {
			log.Printf("error creating session:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, _ := json.Marshal(session)

		_, _ = w.Write(bytes)
	}).Methods("POST").Name("createSession")
	a.HandleFunc("/sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)

		if err != nil {
			log.Printf("error parsing request:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := sessionRep.GetSessionById(r.Context(), id)

		if err != nil {
			log.Printf("error fetching lab:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := r.Context().Value("user").(mysql.User)

		if session.User.ID != user.ID && user.Type != mysql.Professor {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		labUserName := "ubuntu"

		if session.Lab.Type == mysql.Windows {
			labUserName = "administrator"
		}

		bytes, _ := json.Marshal(struct {
			Hostname string `json:"hostname"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Hostname: session.Lab.Hostname,
			Username: labUserName,
			Password: "FlSg5ZJisEecHhMvvBtBPwhjZhdfbnwYjaMR",
		})

		_, _ = w.Write(bytes)
	}).Methods("GET").Name("getSessionInfo")
	a.HandleFunc("/labs/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)

		if err != nil {
			log.Printf("error parsing request:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lab, err := labRep.GetLabById(r.Context(), id)

		if err != nil {
			log.Printf("error fetching lab:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		documentName := "AWS-RunShellScript"
		commands := []string{
			"whoami",
		}
		params := &ssm.SendCommandInput{
			DocumentName: &documentName,
			Parameters: map[string][]string{
				"commands": commands,
			},
			InstanceIds: []string{
				lab.InstanceID,
			},
		}

		out, err := ssmClient.SendCommand(r.Context(), params)

		if err != nil {
			log.Printf("error executing command on lab:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, _ := json.Marshal(out)

		_, _ = w.Write(bytes)
	}).Methods("GET").Name("getStatus")

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
