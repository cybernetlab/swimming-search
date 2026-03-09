package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) GetUserCentres(ctx context.Context, input dto.GetUserCentresInput) (dto.GetCentresOutput, error) {
	centres, err := u.store.GetCentres(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.store.GetCentres: %w", err)
	}
	return collectCentres(input.User.CentreIDs, centres), nil
}
