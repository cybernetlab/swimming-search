package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) AddUserCentre(ctx context.Context, input dto.AddUserCentreInput) (dto.AddUserCentreOutput, error) {
	idx := slices.Index(input.User.CentreIDs, input.CentreID)
	if idx >= 0 {
		return nil, domain.NewErrAlreadyExists("centreID")
	}
	centres, err := u.GetCentres(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.GetCentres: %w", err)
	}
	idx = slices.IndexFunc(centres, func(c domain.Centre) bool { return c.ID == input.CentreID })
	if idx < 0 {
		return nil, domain.ErrNotFound
	}
	input.User.CentreIDs = append(input.User.CentreIDs, input.CentreID)

	err = u.store.PutUser(ctx, *input.User)
	if err != nil {
		return nil, fmt.Errorf("u.store.PutUser: %w", err)
	}
	return collectCentres(input.User.CentreIDs, centres), nil
}
