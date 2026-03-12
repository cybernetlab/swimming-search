package usecase_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/stretchr/testify/require"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func Test_StopSearch_Empty(test *testing.T) {
	t := TestCase{T: test}
	uc := setupStartSearch(&t)
	_, err := uc.StopSearch(t.Context(), dto.StopSearchInput{UserName: "testuser"})
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_StopSearch(test *testing.T) {
	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}

	t := TestCase{T: test}
	uc := setupStartSearch(&t)
	ctx := t.Context()

	err := uc.StartSearch(ctx, dto.StartSearchInput{Search: &search})
	require.NoError(t, err)
	_, err = uc.StopSearch(ctx, dto.StopSearchInput{UserName: "testuser"})
	require.NoError(t, err)
}
