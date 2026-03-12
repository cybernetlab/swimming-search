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

func setupDeleteUser(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	t.Store = new(mocks.Store)
	t.Store.On("GetUser", AnyContext, Any).Return(domain.User{Name: "testuser"}, t.GetUserErr)
	t.Store.On("DeleteUser", AnyContext, Any).Return(t.DeleteUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_DeleteUser_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupDeleteUser(&t)

	input := dto.DeleteUserInput{UserName: "testuser"}
	output, err := uc.DeleteUser(t.Context(), input)
	require.NoError(t, err)
	require.Equal(t, "testuser", output.Name)

	t.Store.AssertExpectations(t)
}

func Test_DeleteUser_NotFound(test *testing.T) {
	t := TestCase{T: test, GetUserErr: domain.ErrNotFound}
	uc := setupDeleteUser(&t)

	input := dto.DeleteUserInput{UserName: "nouser"}
	_, err := uc.DeleteUser(t.Context(), input)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_DeleteUser_DeleteUserError(test *testing.T) {
	t := TestCase{T: test, DeleteUserErr: fmt.Errorf("delete user error")}
	uc := setupDeleteUser(&t)

	input := dto.DeleteUserInput{UserName: "testuser"}
	_, err := uc.DeleteUser(t.Context(), input)
	require.ErrorIs(t, err, t.DeleteUserErr)
}
