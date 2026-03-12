package usecase

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
)

func (u *UseCase) GetUsers(ctx context.Context) (dto.GetUsersOutput, error) {
	output, err := u.store.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("u.store.GetUsers: %w", err)
	}
	slices.SortFunc(output, userSorter)
	return output, nil
}

func userSorter(a, b domain.User) int {
	if a.IsAdmin && !b.IsAdmin {
		return -1
	}
	if !a.IsAdmin && b.IsAdmin {
		return 1
	}
	return cmp.Compare(a.Name, b.Name)
}
