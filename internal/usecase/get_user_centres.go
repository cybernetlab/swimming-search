package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) GetUserCentres(ctx context.Context, input dto.GetUserCentresInput) (dto.GetCentresOutput, error) {
	centres, err := u.GetCentres(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.GetCentres: %w", err)
	}
	return collectCentres(input.User.CentreIDs, centres), nil
}
