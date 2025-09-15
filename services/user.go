package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"chatapp/repository"

	"github.com/charmbracelet/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewUserService(logger *log.Logger, db *sql.DB) *UserService {
	return &UserService{
		DB:     db,
		Logger: logger,
	}
}

type RegisterParams struct {
	Name           string
	Username       string
	Email          string
	Password       string
	VerifyPassword string
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (me *RegisterParams) CleanAndValidate() error {
	me.Name = strings.TrimSpace(me.Name)
	me.Username = strings.TrimSpace(me.Username)
	me.Email = strings.TrimSpace(me.Email)

	return validation.ValidateStruct(me,
		validation.Field(&me.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&me.Username, validation.Required,
			validation.Length(3, 50), validation.Match(usernameRegex).Error("must contain only letters, numbers, or underscore")),
		validation.Field(&me.Email, validation.Required, is.Email),
		validation.Field(&me.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&me.VerifyPassword, validation.Required, validation.By(func(value any) error {
			if v, _ := value.(string); v != me.Password {
				return validation.NewError("validation_password_mismatch", "Passwords do not match")
			}
			return nil
		})),
	)
}

func (me *UserService) Register(params RegisterParams) error {
	if err := params.CleanAndValidate(); err != nil {
		return fmt.Errorf("%w: %w", ErrValidation, err)
	}

	ctx := context.Background()

	tx, err := me.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := repository.New(me.DB).WithTx(tx)

	if ok, err := qtx.CheckUsername(ctx, params.Username); err != nil {
		return err
	} else if ok {
		return validation.Errors{
			"Username": validation.NewError("validation_username_taken", "username is already taken"),
		}
	}

	if ok, err := qtx.CheckEmail(ctx, params.Email); err != nil {
		return err
	} else if ok {
		return validation.Errors{
			"Email": validation.NewError("validation_email_taken", "email is already taken"),
		}
	}

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if err := qtx.InsertUser(ctx, repository.InsertUserParams{
		ID:           uuid.New(),
		Name:         params.Name,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: passwordHash,
	}); err != nil {
		return fmt.Errorf("error insert user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing tx: %w", err)
	}

	return nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

type LoginParams struct {
	Email     string
	Password  string
	UserAgent string
	IpAddress string
}

func (me *LoginParams) CleanAndValidate() error {
	me.Email = strings.TrimSpace(me.Email)
	me.UserAgent = strings.TrimSpace(me.UserAgent)
	me.IpAddress = strings.TrimSpace(me.IpAddress)

	return validation.ValidateStruct(me,
		validation.Field(&me.Email, validation.Required),
		validation.Field(&me.Password, validation.Required),
		validation.Field(&me.UserAgent, validation.Required),
		validation.Field(&me.IpAddress, validation.Required, is.IP),
	)
}

func (me *UserService) Login(params LoginParams) (repository.Session, error) {
	if err := params.CleanAndValidate(); err != nil {
		return repository.Session{}, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	q := repository.New(me.DB)
	ctx := context.Background()

	user, err := q.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.Session{}, ErrUnauthorized
		}
		return repository.Session{}, fmt.Errorf("error getting user: %w", err)
	}

	if !VerifyPassword(params.Password, user.PasswordHash) {
		return repository.Session{}, ErrUnauthorized
	}

	sessionID := uuid.New()
	sessionToken := fmt.Sprintf("%s_%s", sessionID, CreateCryptoRandomHex(32))
	csrfToken := fmt.Sprintf("%s_%s", sessionID, CreateCryptoRandomHex(32))

	session, err := q.InsertSession(ctx, repository.InsertSessionParams{
		ID:           sessionID,
		UserID:       user.ID,
		SessionToken: sessionToken,
		CsrfToken:    csrfToken,
		UserAgent:    params.UserAgent,
		IpAddress:    params.IpAddress,
	})
	if err != nil {
		return repository.Session{}, fmt.Errorf("error inserting session: %w", err)
	}

	return session, nil
}

func CreateCryptoRandomHex(n int) string {
	buffer := make([]byte, n)
	_, _ = rand.Read(buffer)
	return hex.EncodeToString(buffer)
}

func (me *UserService) Authenticate(sessionID uuid.UUID, sessionToken string, csrfToken string) (uuid.UUID, error) {
	q := repository.New(me.DB)
	ctx := context.Background()

	session, err := q.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, ErrUnauthorized
		}
		return uuid.UUID{}, fmt.Errorf("error getting session: %w", err)
	}

	if session.SessionToken != sessionToken || session.CsrfToken != csrfToken {
		return uuid.UUID{}, ErrUnauthorized
	}

	return session.UserID, nil
}

func (me *UserService) UpdateSessionLastActive(sessionID uuid.UUID) error {
	q := repository.New(me.DB)
	ctx := context.Background()
	if err := q.UpdateSessionLastActive(ctx, sessionID); err != nil {
		return fmt.Errorf("error updating session last active: %w", err)
	}
	return nil
}

func (me *UserService) Logout(userID, sessionID uuid.UUID) error {
	q := repository.New(me.DB)
	ctx := context.Background()

	if affectedRows, err := q.DeleteSessionForUser(ctx, repository.DeleteSessionForUserParams{
		ID:     sessionID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	} else if affectedRows == 0 {
		return ErrUnauthorized
	}

	return nil
}

type UpdateUserParams struct {
	UserID         uuid.UUID
	Name           string
	Username       string
	Email          string
	Password       string
	VerifyPassword string
}

func (me *UpdateUserParams) CleanAndValidate() error {
	me.Name = strings.TrimSpace(me.Name)
	me.Username = strings.TrimSpace(me.Username)
	me.Email = strings.TrimSpace(me.Email)

	return validation.ValidateStruct(me,
		validation.Field(&me.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&me.Username, validation.Required,
			validation.Length(3, 50), validation.Match(usernameRegex).Error("must contain only letters, numbers, or underscore")),
		validation.Field(&me.Email, validation.Required, is.Email),
		validation.Field(&me.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&me.VerifyPassword, validation.Required, validation.By(func(value any) error {
			if v, _ := value.(string); v != me.Password {
				return validation.NewError("validation_password_mismatch", "Passwords do not match")
			}
			return nil
		})),
	)
}

func (me *UserService) UpdateUser(params UpdateUserParams) error {
	if err := params.CleanAndValidate(); err != nil {
		return fmt.Errorf("%w: %w", ErrValidation, err)
	}

	ctx := context.Background()

	tx, err := me.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := repository.New(me.DB).WithTx(tx)

	user, err := qtx.GetUserByID(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user not found", ErrNotFound)
		}
		return fmt.Errorf("error getting user by id: %w", err)
	}

	if params.Username != user.Username {
		if ok, err := qtx.CheckUsername(ctx, params.Username); err != nil {
			return err
		} else if ok {
			return validation.Errors{
				"Username": validation.NewError("validation_username_taken", "username is already taken"),
			}
		}
	}

	if params.Email != user.Email {
		if ok, err := qtx.CheckEmail(ctx, params.Email); err != nil {
			return err
		} else if ok {
			return validation.Errors{
				"Email": validation.NewError("validation_email_taken", "email is already taken"),
			}
		}
	}

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if err := qtx.UpdateUser(ctx, repository.UpdateUserParams{
		ID:           params.UserID,
		Name:         params.Name,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: passwordHash,
	}); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (me *UserService) DeleteUser(userID uuid.UUID) error {
	q := repository.New(me.DB)
	ctx := context.Background()

	if err := q.DeleteUserByID(ctx, userID); err != nil {
		return fmt.Errorf("error deleting user by id: %w", err)
	}

	return nil
}
