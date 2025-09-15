package handlers

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	if code == fiber.StatusInternalServerError {
		return c.SendStatus(code)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(code).SendString(err.Error())
}

func WithLogger(logger *slog.Logger) fiber.Handler {
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

func GetCurrentUserID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("UserID").(uuid.UUID)
}

func GetCurrentSessionID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("SessionID").(uuid.UUID)
}
