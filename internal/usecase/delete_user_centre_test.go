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

	"github.com/cybernetlab/swimming-search/test/support"
	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupDeleteUserCentre(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	centres := support.ReadFixture[[]domain.Centre](t.T, "centres")

	t.Store = new(mocks.Store)
	t.Store.On("GetCentres", AnyContext).Return(centres, t.GetCentresErr)
	t.Store.On("PutUser", AnyContext, Any).Return(t.PutUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_DeleteUserCentre_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupDeleteUserCentre(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{3, 10, 5}}
	input := dto.DeleteUserCentreInput{User: &user, CentreID: 10}
	output, err := uc.DeleteUserCentre(t.Context(), input)
	require.NoError(t, err)
	require.Len(t, output, 2)
	require.Equal(t, uint(3), output[0].ID)
	require.Equal(t, uint(5), output[1].ID)

	t.Store.AssertExpectations(t)
}

func Test_DeleteUserCentre_NotFoundCentreID(test *testing.T) {
	t := TestCase{T: test}
	uc := setupDeleteUserCentre(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{10}}
	input := dto.DeleteUserCentreInput{User: &user, CentreID: 3}
	_, err := uc.DeleteUserCentre(t.Context(), input)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_DeleteUserCentre_GetCentresError(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: fmt.Errorf("get centres error")}
	uc := setupDeleteUserCentre(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{10}}
	input := dto.DeleteUserCentreInput{User: &user, CentreID: 10}
	_, err := uc.DeleteUserCentre(t.Context(), input)
	require.ErrorIs(t, err, t.GetCentresErr)
}

func Test_DeleteUserCentre_PutUserError(test *testing.T) {
	t := TestCase{T: test, PutUserErr: fmt.Errorf("put user error")}
	uc := setupDeleteUserCentre(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{10}}
	input := dto.DeleteUserCentreInput{User: &user, CentreID: 10}
	_, err := uc.DeleteUserCentre(t.Context(), input)
	require.ErrorIs(t, err, t.PutUserErr)
}
