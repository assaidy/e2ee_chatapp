package auth

import (
	"chatapp/config"
	"chatapp/repo"
	"chatapp/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type AuthService struct {
	queries *repo.Queries
	logger  *slog.Logger
}

func NewAuthService(logger *slog.Logger, queries *repo.Queries) *AuthService {
	return &AuthService{
		queries: queries,
		logger:  logger,
	}
}

// returns credentials ID and an error if failed
func (me *AuthService) CreateCredentials(params CreateCredentialsParams) (uuid.UUID, error) {
	var zero uuid.UUID
	if err := params.validate(); err != nil {
		return zero, fmt.Errorf("%w: %w", service.ErrValidation, err)
	}

	ctx := context.Background()

	if err := me.queries.Begin(ctx); err != nil {
		return zero, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer me.queries.Rollback(context.Background())

	if ok, err := me.queries.CheckEmail(ctx, params.Email); err != nil {
		return zero, fmt.Errorf("failed to check email: %w", err)
	} else if ok {
		return zero, service.ErrEmailConflict
	}

	passwordHash, err := hashPassword(params.Password)
	if err != nil {
		return zero, fmt.Errorf("failed to hash password: %w", err)
	}

	credentialsID := uuid.New()
	if err := me.queries.InsertCredentials(ctx, repo.InsertCredentialsParams{
		ID:           credentialsID,
		Email:        params.Email,
		PasswordHash: passwordHash,
	}); err != nil {
		return zero, fmt.Errorf("failed to insert credentials: %w", err)
	}

	if err := me.queries.Commit(ctx); err != nil {
		return zero, fmt.Errorf("failed to commit tx: %w", err)
	}

	if err := me.sendVerificationEmail(ctx, params.Email); err != nil {
		me.logger.Error("failed to send verification email", "error", err)
	}

	return credentialsID, nil
}

type CreateCredentialsParams struct {
	Email          string
	Password       string
	VerifyPassword string
}

func (me *CreateCredentialsParams) validate() error {
	return validation.ValidateStruct(me,
		validation.Field(&me.Email, validation.Required, is.Email),
		validation.Field(&me.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&me.VerifyPassword, validation.Required, validation.By(func(value any) error {
			if me.Password != me.VerifyPassword {
				return validation.NewError("validation-password-mismatch", "passwords do not match")
			}
			return nil
		})),
	)
}

func (me *AuthService) sendVerificationEmail(ctx context.Context, email string) error {
	tokenID := uuid.New()
	if err := me.queries.InsertEmailVerificationToken(ctx, repo.InsertEmailVerificationTokenParams{
		ID:        tokenID,
		Email:     email,
		ExpiresAt: time.Now().Add(config.EmailVerificationTokenExpiration),
	}); err != nil {
		return fmt.Errorf("failed to insert verification email token: %w", err)
	}

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", config.AppBaseUrl, tokenID)
	return sendEmail(
		email,
		"Chat App Verification Email",
		fmt.Sprintf(`please <a href="%s"> click here </a> to verify your email.`, verificationLink),
	)
}

func (me *AuthService) StartEmailVerificationCleanupWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(config.EmailVerificationTokenCleanupWorkerTick):
				if err := me.queries.DeleteStaleEmailVerificationTokens(ctx); err != nil {
					me.logger.Error("failed to delete stale email verification tokens", "errors", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (me *AuthService) VerifyEmail(tokenID uuid.UUID) (bool, error) {
	ctx := context.Background()

	token, err := me.queries.GetEmailVerificationTokenByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get email verification token: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		return false, nil
	}

	if err := me.queries.MarkEmailAsVerified(ctx, token.Email); err != nil {
		return false, fmt.Errorf("failed to mark email as verified: %w", err)
	}

	return true, nil
}

func (me *AuthService) Login(email, password string) (repo.Session, error) {
	ctx := context.Background()
	var zero repo.Session

	credentials, err := me.queries.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, service.ErrUnauthorized
		}
		return zero, fmt.Errorf("failed to get credentials by email: %w", err)
	}

	if !verifyPassword(password, credentials.PasswordHash) {
		return zero, service.ErrUnauthorized
	}

	if !credentials.EmailIsVerified {
		return zero, service.ErrEmailNotVerified
	}

	sessionID := uuid.New()
	sessionToken := fmt.Sprintf("%s_%s", sessionID, createRandomHex(32))
	csrfToken := fmt.Sprintf("%s_%s", sessionID, createRandomHex(32))

	session, err := me.queries.InsertSession(ctx, repo.InsertSessionParams{
		ID:            sessionID,
		CredentialsID: credentials.ID,
		Token:         sessionToken,
		CsrfToken:     csrfToken,
	})
	if err != nil {
		return zero, fmt.Errorf("failed to insert session: %w", err)
	}

	return session, nil
}

// returns credentials id
func (me *AuthService) ValidateSession(sessionID uuid.UUID, sessionToken, csrfToken string) (uuid.UUID, error) {
	ctx := context.Background()
	var zero uuid.UUID

	session, err := me.queries.GetSessionByID(ctx, sessionID) 
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, service.ErrUnauthorized
		}
		return zero, fmt.Errorf("failed to get session by id: %w", err)
	}

	return session.CredentialsID, nil
}
