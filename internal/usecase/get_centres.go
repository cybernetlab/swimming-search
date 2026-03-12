package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) GetCentres(ctx context.Context) (dto.GetCentresOutput, error) {
	output, err := u.store.GetCentres(ctx)
	if err == nil {
		return output, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("u.store.GetCentres: %w", err)
	}

	output, err = u.booking.GetCentres(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.booking.GetCentres: %w", err)
	}
	slices.SortFunc(output, centreSorter)

	err = u.store.PutCentres(ctx, output)
	if err != nil {
		return nil, fmt.Errorf("u.store.PutCentres: %w", err)
	}
	return output, nil
}

func centreSorter(a, b domain.Centre) int {
	return int(a.ID) - int(b.ID)
}
