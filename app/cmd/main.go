package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/danutavadanei/nice-lab-go/internal/adapters/mysql"
	"github.com/danutavadanei/nice-lab-go/internal/config"
	"github.com/danutavadanei/nice-lab-go/internal/server"
	"github.com/danutavadanei/nice-lab-go/internal/server/middleware"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// TODO: CHANGE!
const tempPassword = "FlSg5ZJisEecHhMvvBtBPwhjZhdfbnwYjaMR"

func main() {
	v := viper.New()
	v.AutomaticEnv()

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

		user := r.Context().Value("user").(mysql.User)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		out, err := initLabForUser(
			ctx,
			ssmClient,
			&lab,
			&user,
		)

		log.Println(out)

		if err != nil {
			log.Printf("error init lab:  %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := sessionRep.CreateSession(r.Context(), user, lab)

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

		bytes, _ := json.Marshal(struct {
			Hostname string `json:"hostname"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Hostname: session.Lab.Hostname,
			Username: slug.Make(user.Email),
			Password: tempPassword,
		})

		_, _ = w.Write(bytes)
	}).Methods("GET").Name("getSessionInfo")

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

func initLabForUser(ctx context.Context, client *ssm.Client, lab *mysql.Lab, user *mysql.User) (string, error) {
	documentName := "AWS-RunShellScript"
	userName := slug.Make(user.Email)
	commands := []string{
		fmt.Sprintf("adduser --gecos \"\" %s", userName),
		fmt.Sprintf(
			"echo \"%s:%s\" | chpasswd",
			userName,
			tempPassword,
		),
		fmt.Sprintf(
			"/usr/bin/dcv create-session --owner=%s %s",
			userName,
			userName,
		),
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

	sendOut, err := client.SendCommand(ctx, params)

	if err != nil {
		return "", err
	}

	done := make(chan bool)
	var cmdOut *ssm.GetCommandInvocationOutput

	go func() {
	loop:
		cmdOut, err = client.GetCommandInvocation(
			ctx,
			&ssm.GetCommandInvocationInput{
				CommandId:  sendOut.Command.CommandId,
				InstanceId: &lab.InstanceID,
			},
		)

		if err != nil {
			done <- true
		}

		if cmdOut.ResponseCode == -1 {
			time.Sleep(100 * time.Millisecond)
			goto loop
		}

		done <- true
	}()

	select {
	case <-done:
		if err != nil {
			return "", err
		}

		if cmdOut.ResponseCode != 0 {
			return *cmdOut.StandardErrorContent, nil
		}

		return *cmdOut.StandardOutputContent, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
