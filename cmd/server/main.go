package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/FeelsCoderMan/task-manager/internal/api/task"
	_ "github.com/lib/pq"
)

const (
	address = "localhost:8080"
)

type taskLogger struct {
	error *log.Logger
	info  *log.Logger
}

func newTaskLogger() *taskLogger {
	flags := log.Ldate | log.Lshortfile
	errorLogger := log.New(os.Stdout, "ERROR: ", flags)
	infoLogger := log.New(os.Stdout, "INFO: ", flags)

	return &taskLogger{
		error: errorLogger,
		info:  infoLogger,
	}
}

func main() {
	taskLogger := newTaskLogger()

	connStr := "user=postgres password=12345 dbname=task_manager host=localhost sslmode=disable port=5432"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		taskLogger.error.Println("Connection with db is failed: %w", err)
		os.Exit(-1)
	}

	router := http.NewServeMux()
	httpServer := http.Server{
		Addr:    address,
		Handler: router,
	}

	task.RegisterHandlers(router, task.NewService(db))

	taskLogger.info.Printf("Server is running at %s\n", address)
	if err := httpServer.ListenAndServe(); err != nil {
		taskLogger.error.Printf("Could not listen to %s: %s\n", address, err.Error())
		os.Exit(-1)
	}
}
