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

func setupAddUserCentre(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	centres := ReadFixture[[]domain.Centre](t.T, "centres")

	t.Store = new(mocks.Store)
	t.Store.On("GetCentres", AnyContext).Return(centres, t.GetCentresErr)
	t.Store.On("PutUser", AnyContext, Any).Return(t.PutUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_AddUserCentre_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupAddUserCentre(&t)

	user := domain.User{Name: "testuser"}
	input := dto.AddUserCentreInput{User: &user, CentreID: 10}
	output, err := uc.AddUserCentre(t.Context(), input)
	require.NoError(t, err)
	require.Len(t, output, 1)
	require.Equal(t, uint(10), output[0].ID)

	t.Store.AssertExpectations(t)
}

func Test_AddUserCentre_AlreadyExistsCentreID(test *testing.T) {
	t := TestCase{T: test}
	uc := setupAddUserCentre(&t)

	user := domain.User{Name: "testuser", CentreIDs: []uint{10}}
	input := dto.AddUserCentreInput{User: &user, CentreID: 10}
	_, err := uc.AddUserCentre(t.Context(), input)
	var existsErr domain.ErrAlreadyExists
	require.ErrorAs(t, err, &existsErr)
	require.Equal(t, existsErr.Subject, "centreID")
}

func Test_AddUserCentre_InvalidCentreID(test *testing.T) {
	t := TestCase{T: test}
	uc := setupAddUserCentre(&t)

	user := domain.User{Name: "testuser"}
	input := dto.AddUserCentreInput{User: &user, CentreID: 100500}
	_, err := uc.AddUserCentre(t.Context(), input)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_AddUserCentre_GetCentresError(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: fmt.Errorf("get centres error")}
	uc := setupAddUserCentre(&t)

	user := domain.User{Name: "testuser"}
	input := dto.AddUserCentreInput{User: &user, CentreID: 10}
	_, err := uc.AddUserCentre(t.Context(), input)
	require.ErrorIs(t, err, t.GetCentresErr)
}

func Test_AddUserCentre_PutUserError(test *testing.T) {
	t := TestCase{T: test, PutUserErr: fmt.Errorf("put user error")}
	uc := setupAddUserCentre(&t)

	user := domain.User{Name: "testuser"}
	input := dto.AddUserCentreInput{User: &user, CentreID: 10}
	_, err := uc.AddUserCentre(t.Context(), input)
	require.ErrorIs(t, err, t.PutUserErr)
}
