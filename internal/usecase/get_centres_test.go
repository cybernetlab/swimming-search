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

func setupGetCentres(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	centres := ReadFixture[[]domain.Centre](t.T, "centres")

	t.Store = new(mocks.Store)
	t.Store.On("GetCentres", AnyContext).Return(centres, t.GetCentresErr)
	if t.GetCentresErr != nil {
		t.Store.On("PutCentres", AnyContext, Any).Return(t.PutCentresErr)
	}

	t.Booking = new(mocks.Booking)
	if t.GetCentresErr != nil && (t.GetCentresErr == domain.ErrNotFound || t.PutCentresErr != nil) {
		t.Booking.On("GetCentres", AnyContext).Return(centres, t.BookingGetCentresErr)
	}

	return usecase.New(t.Store, t.Booking, nil)
}

func Test_GetCentres_CachedSuccess(test *testing.T) {
	t := TestCase{T: test}
	uc := setupGetCentres(&t)

	output, err := uc.GetCentres(t.Context())
	require.NoError(t, err)
	require.Len(t, output, 15)

	t.Store.AssertExpectations(t)
}

func Test_GetCentres_Success(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: domain.ErrNotFound}
	uc := setupGetCentres(&t)

	output, err := uc.GetCentres(t.Context())
	require.NoError(t, err)
	require.Len(t, output, 15)

	t.Store.AssertExpectations(t)
	t.Booking.AssertExpectations(t)
}

func Test_GetCentres_GetCentresError(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: fmt.Errorf("get centres error")}
	uc := setupGetCentres(&t)

	_, err := uc.GetCentres(t.Context())
	require.ErrorIs(t, err, t.GetCentresErr)
}

func Test_GetCentres_GetCentresBookingError(test *testing.T) {
	t := TestCase{
		T:                    test,
		GetCentresErr:        domain.ErrNotFound,
		BookingGetCentresErr: fmt.Errorf("get centres booking error"),
	}
	uc := setupGetCentres(&t)

	_, err := uc.GetCentres(t.Context())
	require.ErrorIs(t, err, t.BookingGetCentresErr)
}

func Test_GetCentres_PutCentresError(test *testing.T) {
	t := TestCase{T: test, GetCentresErr: domain.ErrNotFound, PutCentresErr: fmt.Errorf("put centres error")}
	uc := setupGetCentres(&t)

	_, err := uc.GetCentres(t.Context())
	require.ErrorIs(t, err, t.PutCentresErr)
}
