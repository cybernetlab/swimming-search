package usecase

import (
	"context"

	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) GetSearch(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, error) {
	for _, search := range u.searches {
		if search.UserName == input.UserName {
			return dto.GetSearchOutput{Search: search}, nil
		}
	}
	return dto.GetSearchOutput{}, domain.ErrNotFound
}
