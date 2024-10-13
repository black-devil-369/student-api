package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/black-devil-369/student-api/internal/config"
	"github.com/black-devil-369/student-api/internal/http/handlers/student"
	"github.com/black-devil-369/student-api/internal/storage/sqllite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// load database
	storage, err := sqllite.New(cfg) // cfg is pointer here

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage is ready", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/student", student.NewUser(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server is starting on ", slog.String("address", cfg.Addr))
	// fmt.Printf("Server is running on %s", cfg.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to Start Server")
		}
	}()

	<-done

	slog.Info("Shutting down the server...")

	// driectly sutdow server
	// by using this - server.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")
}
