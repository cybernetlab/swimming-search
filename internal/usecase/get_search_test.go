package usecase_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"
	"github.com/stretchr/testify/require"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupGetSearch(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	t.Store = new(mocks.Store)
	t.Store.On("PutSearch", AnyContext, Any).Return(nil)
	t.Store.On("DeleteSearch", AnyContext, Any).Return(nil).Maybe()

	booking := new(mocks.Booking)
	booking.On("StartSearchCourses", AnyContext, Any, Any).Return()

	return usecase.New(t.Store, booking, nil)
}

func Test_GetSearch_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.NoError(t, err)

	s, err := uc.GetSearch(t.Context(), dto.GetSearchInput{UserName: "testuser"})
	require.NoError(t, err)
	require.Equal(t, "test query", s.NameQuery)

	_, err = uc.GetSearch(t.Context(), dto.GetSearchInput{UserName: "otheruser"})
	require.ErrorIs(t, err, domain.ErrNotFound)

	t.Store.AssertExpectations(t)
}
