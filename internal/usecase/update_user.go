package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	_, err := u.store.GetUser(ctx, input.User.Name)
	if err != nil {
		return fmt.Errorf("u.store.GetUser: %w", err)
	}
	err = u.store.PutUser(ctx, *input.User)
	if err != nil {
		return fmt.Errorf("u.store.PutUser: %w", err)
	}
	return nil
}
