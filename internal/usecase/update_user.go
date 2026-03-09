package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	err := u.store.PutUser(ctx, *input.User)
	if err != nil {
		return fmt.Errorf("u.store.PutUser: %w", err)
	}
	return nil
}
