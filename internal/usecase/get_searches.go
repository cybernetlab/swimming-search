package usecase

import (
	"context"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
)

func (u *UseCase) GetSearches(ctx context.Context) (dto.GetSearchesOutput, error) {
	node, err := domain.ContextNodeID(ctx)
	if err != nil {
		return nil, fmt.Errorf("domain.ContextNodeID: %w", err)
	}
	output, err := u.store.GetSearches(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("u.store.GetSearches: %w", err)
	}
	return output, nil
}
