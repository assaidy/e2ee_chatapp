package handler

import (
	"chatapp/service"
	"chatapp/service/auth"
	"chatapp/service/user"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *auth.AuthService
	userService *user.UserService
}

func NewAuthHandler(authService *auth.AuthService, userService *user.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (me *AuthHandler) HandleRegister(c *fiber.Ctx) error {
	var (
		name           = strings.TrimSpace(c.FormValue("name"))
		username       = strings.TrimSpace(c.FormValue("username"))
		email          = strings.TrimSpace(c.FormValue("email"))
		password       = c.FormValue("password")
		verifyPassword = c.FormValue("verify-password")
	)

	credentialsID, err := me.authService.CreateCredentials(auth.CreateCredentialsParams{
		Email:          email,
		Password:       password,
		VerifyPassword: verifyPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			if errs, ok := service.ExtractValidationErrorsMap(err); ok {
				return c.Status(fiber.StatusBadRequest).JSON(errs)
			}
			return fmt.Errorf("failed to exctract validation errors")
		case errors.Is(err, service.ErrEmailConflict):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"email": "email already exists",
			})
		}
		return fmt.Errorf("failed to register user: %w", err)
	}

	if err := me.userService.CreateUser(user.CreateUserParams{
		Name:          name,
		Username:      username,
		CredentialsID: credentialsID,
	}); err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			if errs, ok := service.ExtractValidationErrorsMap(err); ok {
				return c.Status(fiber.StatusBadRequest).JSON(errs)
			}
			return fmt.Errorf("failed to exctract validation errors")
		case errors.Is(err, service.ErrUsernameConflict):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"username": "username already exists",
			})
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (me *AuthHandler) HandleVerifyEmail(c *fiber.Ctx) error {
	tokenQuery := c.Query("token")
	tokenID, err := uuid.Parse(tokenQuery)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid token")
	}

	if ok, err := me.authService.VerifyEmail(tokenID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	} else if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("invalid or expired token")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (me *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var (
		email    = strings.TrimSpace(c.FormValue("email"))
		password = c.FormValue("password")
	)

	session, err := me.authService.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			return fiber.ErrUnauthorized
		case errors.Is(err, service.ErrEmailNotVerified):
			return c.Status(fiber.StatusForbidden).SendString("email is not verified")
		}
		return fmt.Errorf("failed to login: %w", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session-id",
		Value:    session.ID.String(),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "session-token",
		Value:    session.Token,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:  "csrf-token",
		Value: session.CsrfToken,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (me *AuthHandler) WithSession(c *fiber.Ctx) error {
	var (
		sessionIDstr = c.Cookies("session-id")
		sessionToken = c.Cookies("session-token")
		csrfToken    = c.Get("X-CSRF-Token")
	)

	sessionID, err := uuid.Parse(sessionIDstr)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if sessionToken == "" || csrfToken == "" {
		return fiber.ErrUnauthorized
	}

	credentialsID, err := me.authService.ValidateSession(sessionID, sessionToken, csrfToken)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("failed to validate sessoin: %w", err)
	}

	c.Locals("auth.credentialsID", credentialsID)
	return c.Next()
}

func getCurrentUserCredentialsID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("auth.credentialsID").(uuid.UUID)
}
