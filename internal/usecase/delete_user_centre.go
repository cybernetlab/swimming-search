package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) DeleteUserCentre(ctx context.Context, input dto.DeleteUserCentreInput) (dto.DeleteUserCentreOutput, error) {
	idx := slices.Index(input.User.CentreIDs, input.CentreID)
	if idx < 0 {
		return nil, domain.ErrNotFound
	}
	input.User.CentreIDs = slices.Delete(input.User.CentreIDs, idx, idx+1)
	centres, err := u.GetCentres(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.GetCentres: %w", err)
	}
	err = u.store.PutUser(ctx, *input.User)
	if err != nil {
		return nil, fmt.Errorf("u.store.PutUser: %w", err)
	}
	return collectCentres(input.User.CentreIDs, centres), nil
}
