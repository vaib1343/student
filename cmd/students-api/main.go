package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vaib1343/students-api/internal/config"
	"github.com/vaib1343/students-api/internal/http/handlers/student"
	"github.com/vaib1343/students-api/internal/storage/sqlite"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// connect to database
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// listen server

	route := http.NewServeMux()

	route.HandleFunc("POST /api/students", student.New(storage))
	route.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	route.HandleFunc("GET /api/students", student.GetStudents(storage))

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: route,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("server is running on port %s\n", cfg.HTTPServer.Addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("failed to start server: %v", err)
		}
	}()

	<-done
	slog.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	slog.Info("server is shut down")

	// start server
}
