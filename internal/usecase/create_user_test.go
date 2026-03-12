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

func setupCreateUser(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	t.Store = new(mocks.Store)
	t.Store.On("GetUser", AnyContext, Any).Return(domain.User{Name: "testuser"}, t.GetUserErr)
	t.Store.On("PutUser", AnyContext, Any).Return(t.PutUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_CreateUser_Success(test *testing.T) {
	t := TestCase{T: test, GetUserErr: domain.ErrNotFound}
	uc := setupCreateUser(&t)

	input := dto.CreateUserInput{UserName: "newuser"}
	output, err := uc.CreateUser(t.Context(), input)
	require.NoError(t, err)
	require.Equal(t, "newuser", output.Name)
	require.Equal(t, false, output.IsAdmin)

	t.Store.AssertExpectations(t)
}

func Test_CreateUser_AlreadyExists(test *testing.T) {
	t := TestCase{T: test}
	uc := setupCreateUser(&t)

	input := dto.CreateUserInput{UserName: "testuser"}
	_, err := uc.CreateUser(t.Context(), input)
	var existsErr domain.ErrAlreadyExists
	require.ErrorAs(t, err, &existsErr)
	require.Equal(t, existsErr.Subject, "user")
}

func Test_CreateUser_GetUserError(test *testing.T) {
	t := TestCase{T: test, GetUserErr: fmt.Errorf("get user error")}
	uc := setupCreateUser(&t)

	input := dto.CreateUserInput{UserName: "newuser"}
	_, err := uc.CreateUser(t.Context(), input)
	require.ErrorIs(t, err, t.GetUserErr)
}

func Test_CreateUser_PutUserError(test *testing.T) {
	t := TestCase{T: test, GetUserErr: domain.ErrNotFound, PutUserErr: fmt.Errorf("put user error")}
	uc := setupCreateUser(&t)

	input := dto.CreateUserInput{UserName: "newuser"}
	_, err := uc.CreateUser(t.Context(), input)
	require.ErrorIs(t, err, t.PutUserErr)
}
