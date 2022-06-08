package main

import (
	"context"
	"github.com/danutavadanei/localstack-go-playground/internal/adapters/aws"
	"github.com/danutavadanei/localstack-go-playground/internal/config"
	"github.com/danutavadanei/localstack-go-playground/internal/server"
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
	v.AutomaticEnv()
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	cfg := config.NewAppConfig(v)

	awsClient := aws.NewClient(cfg.AWSConfig)

	m := mux.NewRouter()

	m.HandleFunc("/s3/buckets", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := awsClient.ListBuckets(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("GET").Name("listBuckets")

	m.HandleFunc("/s3/buckets/{bucket}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]

		bytes, err := awsClient.ListBucketObjects(r.Context(), &bucket)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("GET").Name("listBucketObjects")

	m.HandleFunc("/s3/buckets/{bucket}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]

		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB

		if err != nil {
			log.Printf("error parsing request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f, h, err := r.FormFile("file")

		if err != nil {
			log.Printf("error parsing request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, err := awsClient.UploadFileToBucket(r.Context(), &bucket, &h.Filename, f)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}).Methods("POST", "PUT").Name("putBucketObject")

	m.HandleFunc("/s3/buckets/{bucket}/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket, key := vars["bucket"], vars["key"]

		bytesWritten, err := awsClient.SinkFileToWriter(r.Context(), &bucket, &key, w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Download of \"%s\" complete. Wrote %s bytes", key, strconv.FormatInt(bytesWritten, 10))
	}).Methods("GET").Name("getBucketObject")

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
