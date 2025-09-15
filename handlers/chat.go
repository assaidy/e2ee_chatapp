package handlers

import (
	"log/slog"

	"chatapp/services"
)

type ChatHandler struct {
	Logger      *slog.Logger
	ChatService *services.ChatService
}

func NewChatHandler(logger *slog.Logger, chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		Logger:      logger,
		ChatService: chatService,
	}
}
