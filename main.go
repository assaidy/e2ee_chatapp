package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatapp/env"
	"chatapp/handlers"
	"chatapp/services"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		Formatter:       log.TextFormatter,
		ReportTimestamp: true,
	})

	db, err := sql.Open("postgres", env.PgUrl)
	if err != nil {
		logger.Fatal("error connecting to postgres db", "err", err)
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		logger.Fatal("error pinging postgres db", "err", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetConnMaxIdleTime(15 * time.Minute)

	userService := services.NewUserService(logger, db)
	userHandler := handlers.NewUserHandler(logger, userService)

	server := fiber.New(fiber.Config{
		ErrorHandler: handlers.FiberErrorHandler,
		Prefork:      false,
	})
	server.Use(handlers.WithLogger(slog.New(logger)))

	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.Post("users/register", userHandler.HandleRegister)
			v1.Post("users/login", userHandler.HandleLogin)
			v1.Post("users/logout", userHandler.WithAuthentication, userHandler.HandleLogout)
			v1.Put("/users", userHandler.WithAuthentication, userHandler.HandleUpdateUser)
			v1.Delete("/users", userHandler.WithAuthentication, userHandler.HandleDeleteUser)
		}
	}

	// ===========================================================================================
	// start server
	// ===========================================================================================
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	listenErrChan := make(chan error, 1)
	go func() {
		serverAddr := fmt.Sprintf("localhost:%d", env.ServerPort)
		logger.Info("starting server", "addr", serverAddr)
		if err := server.Listen(serverAddr); err != nil {
			listenErrChan <- err
		}
	}()

	select {
	case err := <-listenErrChan:
		logger.Fatal("error running app", "error", err)
	case <-exitChan:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Fatal("failed to shutdown server", "error", err)
	}
}
