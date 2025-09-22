package handler

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func WithLogging(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		logger.Info("request",
			"took", duration,
			"method", method,
			"path", path,
			"status", status,
			"ip", ip,
			"error", err,
		)

		return err
	}
}
