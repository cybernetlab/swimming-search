package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) CreateUser(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error) {
	_, err := u.store.GetUser(ctx, input.UserName)
	if !errors.Is(err, domain.ErrNotFound) {
		if err == nil {
			return dto.CreateUserOutput{}, domain.NewErrAlreadyExists("user")
		}
		return dto.CreateUserOutput{}, fmt.Errorf("u.store.GetUser: %w", err)
	}
	user := domain.User{
		Name:    input.UserName,
		IsAdmin: input.IsAdmin,
	}
	err = u.store.PutUser(ctx, user)
	if err != nil {
		return dto.CreateUserOutput{}, fmt.Errorf("u.store.PutUser: %w", err)
	}
	return dto.CreateUserOutput{User: user}, nil
}
