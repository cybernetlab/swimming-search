package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) DeleteUser(ctx context.Context, input dto.DeleteUserInput) (dto.DeleteUserOutput, error) {
	user, err := u.store.GetUser(ctx, input.UserName)
	if err != nil {
		return dto.DeleteUserOutput{}, fmt.Errorf("u.store.GetUser: %w", err)
	}
	err = u.store.DeleteUser(ctx, input.UserName)
	if err != nil {
		return dto.DeleteUserOutput{}, fmt.Errorf("u.store.DeleteUser: %w", err)
	}
	return dto.DeleteUserOutput{User: user}, nil
}
