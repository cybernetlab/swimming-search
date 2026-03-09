package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) DeleteUser(ctx context.Context, input dto.GetUserInput) (dto.GetUserOutput, error) {
	user, err := u.store.GetUser(ctx, input.UserName)
	if err != nil {
		return dto.GetUserOutput{}, fmt.Errorf("u.store.GetUser: %w", err)
	}
	err = u.store.DeleteUser(ctx, input.UserName)
	if err != nil {
		return dto.GetUserOutput{}, fmt.Errorf("u.store.DeleteUser: %w", err)
	}
	return dto.GetUserOutput{User: user}, nil
}
