package services

import (
	"database/sql"
	"log/slog"
)

type ChatService struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewChatService(logger *slog.Logger, db *sql.DB) *ChatService {
	return &ChatService{
		DB:     db,
		Logger: logger,
	}
}
