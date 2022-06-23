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

	a := m.PathPrefix("/").Subrouter()
	a.Use(authMiddleware.Middleware)

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

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
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
			Username: user.UserName,
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
	if lab.Type == mysql.Windows {
		return initLabForUserWindows(ctx, client, lab, user)
	}

	return initLabForUserLinux(ctx, client, lab, user)
}

func initLabForUserWindows(ctx context.Context, client *ssm.Client, lab *mysql.Lab, user *mysql.User) (string, error) {
	documentName := "AWS-RunPowerShellScript"

	commands := []string{
		fmt.Sprintf("New-LocalUser -Name \"%[1]s\" -NoPassword -FullName \"%[1]s\"", user.UserName),
		fmt.Sprintf("net user \"%s\" \"%s\"", user.UserName, tempPassword),
		fmt.Sprintf(
			".\"C:\\Program Files\\NICE\\DCV\\Server\\bin\\dcv.exe\" create-session --owner=%[1]s %[1]s",
			user.UserName,
		),
		fmt.Sprintf("md \"Z:\\%s\" 2>NUL", user.UserName),
		fmt.Sprintf("md \"Z:\\%s\\windows\" 2>NUL", user.UserName),
		fmt.Sprintf("$shortcut=(New-Object -ComObject WScript.Shell).CreateShortcut('C:\\Users\\%[1]s\\Desktop\\DCV-Storage.lnk');$shortcut.TargetPath='Z:\\%[1]s\\Windows';$shortcut.Save()", user.UserName),
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

		if cmdOut != nil && cmdOut.ResponseCode == -1 {
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

func initLabForUserLinux(ctx context.Context, client *ssm.Client, lab *mysql.Lab, user *mysql.User) (string, error) {
	documentName := "AWS-RunShellScript"
	commands := []string{
		fmt.Sprintf("adduser --gecos \"\" %s", user.UserName),
		fmt.Sprintf(
			"echo \"%s:%s\" | chpasswd",
			user.UserName,
			tempPassword,
		),
		fmt.Sprintf("/usr/bin/dcv create-session --owner=%[1]s %[1]s", user.UserName),
		fmt.Sprintf("mkdir -p /var/fsx/%s/linux", user.UserName),
		fmt.Sprintf("mkdir -p /home/%s/Desktop", user.UserName),
		fmt.Sprintf("chown -R %[1]s:%[1]s /home/%[1]s/Desktop", user.UserName),
		fmt.Sprintf("ln -s /var/fsx/%[1]s/linux /home/%[1]s/Desktop/NiceLabData", user.UserName),
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

	time.Sleep(100 * time.Millisecond)
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

		if cmdOut != nil && cmdOut.ResponseCode == -1 {
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
