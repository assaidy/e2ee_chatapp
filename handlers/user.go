package handlers

import (
	"errors"
	"fmt"

	"chatapp/services"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Logger      *log.Logger
	UserService *services.UserService
}

func NewUserHandler(logger *log.Logger, userService *services.UserService) *UserHandler {
	return &UserHandler{
		Logger:      logger,
		UserService: userService,
	}
}

func (me *UserHandler) HandleRegister(c *fiber.Ctx) error {
	params := services.RegisterParams{
		Name:           c.FormValue("name"),
		Username:       c.FormValue("username"),
		Email:          c.FormValue("email"),
		Password:       c.FormValue("password"),
		VerifyPassword: c.FormValue("verify_password"),
	}

	if err := me.UserService.Register(params); err != nil {
		if errors.Is(err, services.ErrValidation) {
			if m := services.ExtractValidationErrors(err); m != nil {
				return c.Status(fiber.StatusBadRequest).JSON(m)
			} else {
				return fmt.Errorf("couldn't extract validation errors")
			}
		}
		return err
	}

	// TODO: verify email

	return c.SendStatus(fiber.StatusCreated)
}

func (me *UserHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	session, err := me.UserService.Login(services.LoginParams{
		Email:     email,
		Password:  password,
		UserAgent: string(c.Request().Header.UserAgent()),
		IpAddress: c.IP(),
	})
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    session.CsrfToken,
		HTTPOnly: false,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (me *UserHandler) WithAuthentication(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Cookies("session_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	sessionToken := c.Cookies("session_token")
	csrfToken := c.Get("X-CSRF-Token")
	if sessionToken == "" || csrfToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, err := me.UserService.Authenticate(sessionID, sessionToken, csrfToken)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	if err := me.UserService.UpdateSessionLastActive(sessionID); err != nil {
		return err
	}

	c.Locals("SessionID", sessionID)
	c.Locals("UserID", userID)
	return c.Next()
}

func (me *UserHandler) HandleLogout(c *fiber.Ctx) error {
	userID := GetCurrentUserID(c)
	sessionID := GetCurrentSessionID(c)

	if query := c.Query("session_id"); query != "" {
		// logout from a specific session for the current user
		if id, err := uuid.Parse(query); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid session id")
		} else {
			sessionID = id
		}
	} else {
		// logout from the current session
		c.ClearCookie("session_id")
		c.ClearCookie("session_token")
		c.ClearCookie("csrf_token")
	}

	if err := me.UserService.Logout(userID, sessionID); err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (me *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userID := GetCurrentUserID(c)
	params := services.UpdateUserParams{
		UserID:         userID,
		Name:           c.FormValue("name"),
		Username:       c.FormValue("username"),
		Email:          c.FormValue("email"),
		Password:       c.FormValue("password"),
		VerifyPassword: c.FormValue("verify_password"),
	}

	if err := me.UserService.UpdateUser(params); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.ClearCookie("session_id")
			c.ClearCookie("session_token")
			c.ClearCookie("csrf_token")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (me *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := GetCurrentUserID(c)
	if err := me.UserService.DeleteUser(userID); err != nil {
		return err
	}

	c.ClearCookie("session_id")
	c.ClearCookie("session_token")
	c.ClearCookie("csrf_token")
	return c.SendStatus(fiber.StatusOK)
}
