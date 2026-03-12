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

func setupGetUser(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	users := ReadFixture[[]domain.User](t.T, "users")
	t.Store = new(mocks.Store)
	t.Store.On("GetUser", AnyContext, Any).Return(users[1], t.GetUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_GetUser_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetUser(&t)

	output, err := uc.GetUser(t.Context(), dto.GetUserInput{UserName: "testuser1"})
	require.NoError(t, err)
	require.Equal(t, "testuser1", output.Name)
}

func Test_GetUser_GetUserError(test *testing.T) {
	t := TestCase{T: test, GetUserErr: fmt.Errorf("get user error")}
	uc := setupGetUser(&t)

	_, err := uc.GetUser(t.Context(), dto.GetUserInput{UserName: "testuser1"})
	require.ErrorIs(t, err, t.GetUserErr)
}
