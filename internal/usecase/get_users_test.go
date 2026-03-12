package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupGetUsers(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()
	ctx := t.Context()

	users := ReadFixture[[]domain.User](t.T, "users")
	t.Store = new(mocks.Store)
	t.Store.On("GetUsers", ctx).Return(users, t.GetUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_GetUsers_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetUsers(&t)

	output, err := uc.GetUsers(t.Context())
	require.NoError(t, err)
	require.Len(t, output, 4)
	require.Equal(t, "testuser2", output[0].Name)
	require.Equal(t, "testuser3", output[1].Name)
	require.Equal(t, "testuser1", output[2].Name)
	require.Equal(t, "testuser4", output[3].Name)
}

func Test_GetUsers_GetUserError(test *testing.T) {
	t := TestCase{T: test, GetUserErr: fmt.Errorf("get users error")}
	uc := setupGetUsers(&t)

	_, err := uc.GetUsers(t.Context())
	require.ErrorIs(t, err, t.GetUserErr)
}
