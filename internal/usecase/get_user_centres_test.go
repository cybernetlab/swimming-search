package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupGetUserCentres(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	centres := ReadFixture[[]domain.Centre](t.T, "centres")

	t.Store = new(mocks.Store)
	t.Store.On("GetCentres", AnyContext).Return(centres, t.GetCentresErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_GetUserCentres_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetUserCentres(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{3, 10, 100}}
	output, err := uc.GetUserCentres(t.Context(), dto.GetUserCentresInput{User: &user})
	require.NoError(t, err)
	require.Len(t, output, 2)
	require.Equal(t, "Scott Hall", output[0].Name)
	require.Equal(t, "Kirkstall", output[1].Name)
}

func Test_GetUserCentres_GetCentresError(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: fmt.Errorf("get centres error")}
	uc := setupGetUserCentres(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{3, 10}}
	_, err := uc.GetUserCentres(t.Context(), dto.GetUserCentresInput{User: &user})
	require.ErrorIs(t, err, t.GetCentresErr)
}
