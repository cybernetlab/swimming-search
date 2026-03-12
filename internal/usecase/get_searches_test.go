package usecase_test

import (
	"fmt"
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"
	"github.com/stretchr/testify/require"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupGetSearches(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	searches := ReadFixture[[]domain.Search](t.T, "searches")
	t.Store = new(mocks.Store)
	t.Store.On("GetSearches", AnyContext, Any).Return(searches, t.GetSearchesErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_GetSearches_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetSearches(&t)

	s, err := uc.GetSearches(domain.WithNodeID(t.Context(), "node1"))
	require.NoError(t, err)
	require.Len(t, s, 2)
	require.Equal(t, "search-1", s[0].NameQuery)
}

func Test_GetSearches_InvalidContext(test *testing.T) {
	var ctxErr domain.ErrInvalidContext
	t := TestCase{T: test}
	uc := setupGetSearches(&t)

	_, err := uc.GetSearches(t.Context())
	require.ErrorAs(t, err, &ctxErr)
	require.Equal(t, "nodeID", ctxErr.Field)
}

func Test_GetSearches_GetSearchesErr(test *testing.T) {
	t := TestCase{T: test, GetSearchesErr: fmt.Errorf("get searches error")}
	uc := setupGetSearches(&t)

	_, err := uc.GetSearches(domain.WithNodeID(t.Context(), "node1"))
	require.ErrorIs(t, err, t.GetSearchesErr)
}
