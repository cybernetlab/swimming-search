package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/mocks"
	"github.com/cybernetlab/swimming-search/pkg/otel"
	"github.com/stretchr/testify/require"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func setupStartSearch(t *TestCase) *usecase.UseCase {
	otel.SilentModeInit()

	t.Store = new(mocks.Store)
	t.Store.On("PutSearch", AnyContext, Any).Return(t.PutSearchErr)
	t.Store.On("DeleteSearch", AnyContext, Any).Return(t.DeleteSearchError)

	t.Booking = new(mocks.Booking)
	t.Booking.On("StartSearchCourses", AnyContext, Any, Any).Return().Run(func(args mock.Arguments) {
		t.Courses = args.Get(2).(chan<- domain.Course)
	})

	t.Bot = new(mocks.Bot)
	t.Bot.On("Send", AnyContext, Any).Return(t.BotSendErr)

	return usecase.New(t.Store, t.Booking, t.Bot)
}

func Test_StartSearch_Success(test *testing.T) {
	t := TestCase{T: test}
	uc := setupStartSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.NoError(t, err)
}

func Test_StartSearch_AlreadyExists(test *testing.T) {
	var existsErr domain.ErrAlreadyExists
	t := TestCase{T: test}
	uc := setupStartSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.NoError(t, err)

	err = uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.ErrorAs(t, err, &existsErr)
	require.Equal(t, existsErr.Subject, "search")
}

func Test_StartSearch_NoCentreIDs(test *testing.T) {
	t := TestCase{T: test}
	uc := setupStartSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query"}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.ErrorIs(t, err, domain.ErrEmptyCentreIDs)
}

func Test_StartSearch_PutError(test *testing.T) {
	t := TestCase{T: test, PutSearchErr: fmt.Errorf("put error")}
	uc := setupStartSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.ErrorIs(t, err, t.PutSearchErr)
}

func Test_StartSearch_DeleteError(test *testing.T) {
	t := TestCase{T: test, DeleteSearchError: fmt.Errorf("delete error")}
	uc := setupStartSearch(&t)

	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}
	err := uc.StartSearch(t.Context(), dto.StartSearchInput{Search: &search})
	require.NoError(t, err)
}

func Test_StartSearch_ChanClose(test *testing.T) {
	sched := domain.CourseSchedule{Type: "dayOfWeek", Time: struct {
		Start string
		End   string
	}{Start: "15:00:00", End: "16:00:00"}}
	search := domain.Search{UserName: "testuser", NameQuery: "test query", CentreIDs: []uint{10}}

	t := TestCase{T: test, DeleteSearchError: fmt.Errorf("delete error"), BotSendErr: fmt.Errorf("send error")}
	uc := setupStartSearch(&t)
	ctx := t.Context()

	err := uc.StartSearch(ctx, dto.StartSearchInput{Search: &search})
	require.NoError(t, err)
	t.Courses <- domain.Course{Name: "course1", Schedule: sched}
	close(t.Courses)
	_, err = uc.StopSearch(ctx, dto.StopSearchInput{UserName: "testuser"})
	require.NoError(t, err)

	time.Sleep(100 * time.Microsecond)

	err = uc.StartSearch(ctx, dto.StartSearchInput{Search: &search})
	require.NoError(t, err)
	t.Courses <- domain.Course{Name: "course2"}
	close(t.Courses)
	_, err = uc.StopSearch(ctx, dto.StopSearchInput{UserName: "testuser"})
	require.NoError(t, err)

	t.Store.AssertExpectations(t)
}
