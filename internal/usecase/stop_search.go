package usecase

import (
	"context"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) StopSearch(ctx context.Context, input dto.StopSearchInput) (dto.StopSearchOutput, error) {
	for _, search := range u.searches {
		if search.UserName == input.UserName {
			if search.Cancel != nil {
				search.Cancel()
			}
			return dto.StopSearchOutput{Search: search}, nil
		}
	}
	return dto.StopSearchOutput{}, domain.ErrNotFound
}
