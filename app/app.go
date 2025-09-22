package app

import (
	"chatapp/config"
	"chatapp/handler"
	"chatapp/service/auth"
	"chatapp/service/user"

	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	logger      *slog.Logger
	authService *auth.AuthService
	userService *user.UserService
}

func NewApp(logger *slog.Logger, authService *auth.AuthService, userService *user.UserService) *App {
	return &App{
		logger:      logger,
		authService: authService,
		userService: userService,
	}
}

func (me *App) Run() error {
	me.logger.Info("loading routes")
	server := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
		Prefork:      false,
	})
	server.Use(handler.WithLogging(me.logger))
	me.loadAuthRoutes(server)
	me.loadUserRoutes(server)

	listenErrChan := make(chan error, 1)
	go func() {
		me.logger.Info("starting server", "addr", config.ListenAddr)
		if err := server.Listen(config.ListenAddr); err != nil {
			listenErrChan <- err
		}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErrChan:
		return fmt.Errorf("listen error: %w", err)

	case <-exitChan:
		me.logger.Info("starting server shutdown")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
	}

	return nil
}

func (me *App) loadAuthRoutes(server *fiber.App) {
	ah := handler.NewAuthHandler(me.authService, me.userService)

	server.Post("/register", ah.HandleRegister)
	server.Get("/verify-email", ah.HandleVerifyEmail)
	server.Get("/login", ah.HandleLogin)
}

// TODO:
func (me *App) loadUserRoutes(server *fiber.App) {}
