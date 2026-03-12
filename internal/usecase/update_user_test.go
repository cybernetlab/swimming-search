package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupUpdateUser(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	t.Store = new(mocks.Store)
	t.Store.On("GetUser", AnyContext, Any).Return(domain.User{Name: "testuser"}, t.GetUserErr)

	return usecase.New(t.Store, nil, nil)
}

func Test_UpdateUser_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupUpdateUser(&t)

	input := dto.UpdateUserInput{User: &domain.User{Name: "testuser", IsAdmin: true}}
	t.Store.On("PutUser", AnyContext, Any).Return(nil).Run(func(args mock.Arguments) {
		require.Equal(t, true, args.Get(1).(domain.User).IsAdmin)
	})
	err := uc.UpdateUser(t.Context(), input)
	require.NoError(t, err)

	t.Store.AssertExpectations(t)
}

func Test_UpdateUser_NotFound(test *testing.T) {
	t := TestCase{T: test, GetUserErr: domain.ErrNotFound}
	uc := setupUpdateUser(&t)

	input := dto.UpdateUserInput{User: &domain.User{Name: "testuser"}}
	err := uc.UpdateUser(t.Context(), input)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_UpdateUser_PutUserError(test *testing.T) {
	putErr := fmt.Errorf("delete user error")
	t := TestCase{T: test}
	uc := setupUpdateUser(&t)
	t.Store.On("PutUser", AnyContext, Any).Return(putErr)

	input := dto.UpdateUserInput{User: &domain.User{Name: "testuser"}}
	err := uc.UpdateUser(t.Context(), input)
	require.ErrorIs(t, err, putErr)
}
