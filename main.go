package main

import (
	"chatapp/app"
	"chatapp/db"
	"chatapp/repo"
	"chatapp/service/auth"
	"chatapp/service/user"
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	logHander := log.NewWithOptions(os.Stdout, log.Options{
		Formatter:       log.TextFormatter,
		ReportTimestamp: true,
	})
	logger := slog.New(logHander)

	workersCtx, workersCancel := context.WithCancel(context.Background())
	defer workersCancel()

	authService := auth.NewAuthService(logger, repo.New(db.DB))
	authService.StartEmailVerificationCleanupWorker(workersCtx)

	userService := user.NewUserService(repo.New(db.DB))

	app := app.NewApp(
		logger,
		authService,
		userService,
	)
	if err := app.Run(); err != nil {
		logger.Error("failed to run app", "error", err)
	}
}
