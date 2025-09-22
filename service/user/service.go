package user

import (
	"chatapp/repo"
	"chatapp/service"
	"context"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type UserService struct {
	queries *repo.Queries
}

func NewUserService(queries *repo.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

func (me *UserService) CreateUser(params CreateUserParams) error {
	if err := params.validate(); err != nil {
		return fmt.Errorf("%w: %w", service.ErrValidation, err)
	}

	ctx := context.Background()

	if err := me.queries.Begin(ctx); err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer me.queries.Rollback(context.Background())

	if ok, err := me.queries.CheckUsername(ctx, params.Username); err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	} else if ok {
		return service.ErrUsernameConflict
	}

	if err := me.queries.InsertUser(ctx, repo.InsertUserParams{
		ID:            uuid.New(),
		Name:          params.Name,
		Username:      params.Username,
		CredentialsID: params.CredentialsID,
	}); err != nil {
		return fmt.Errorf("failed to insert profile: %w", err)
	}

	if err := me.queries.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

type CreateUserParams struct {
	Name          string
	Username      string
	CredentialsID uuid.UUID
}

var usernameRegex = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func (me *CreateUserParams) validate() error {
	return validation.ValidateStruct(me,
		validation.Field(&me.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&me.Username, validation.Required, validation.Length(2, 50),
			validation.Match(usernameRegex).Error("must contain only letters, numbers, or underscore")),
		validation.Field(&me.CredentialsID, validation.Required, is.UUID),
	)
}
