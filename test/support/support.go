package support

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/mocks"
)

const Any = mock.Anything

var AnyContext = mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })

type TestCase struct {
	*testing.T
	BotSendErr           error
	BookingGetCentresErr error
	DeleteSearchError    error
	DeleteUserErr        error
	GetCentresErr        error
	GetSearchesErr       error
	GetUserErr           error
	PutUserErr           error
	PutCentresErr        error
	PutSearchErr         error
	Store                *mocks.Store
	Booking              *mocks.Booking
	Bot                  *mocks.Bot
	Courses              chan<- domain.Course
}

type Subst struct {
	Key   string
	Value string
}

func FixturePath(t *testing.T, fixture string) string {
	_, b, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	baseDir := filepath.Join(filepath.Dir(b), "..", "fixtures")
	return filepath.Join(baseDir, fixture)
}

func GetFixture(t *testing.T, fixture string, subst ...Subst) []byte {
	if !filepath.IsAbs(fixture) {
		fixture = FixturePath(t, fixture)
	}
	payload, err := os.ReadFile(fixture)
	assert.NoError(t, err)
	result := string(payload)
	for _, s := range subst {
		result = strings.ReplaceAll(result, "<"+s.Key+">", s.Value)
	}
	return []byte(result)
}

func ReadFixture[T domain.Course | []domain.Centre | []domain.Search | []domain.User](
	t *testing.T,
	fixture string,
	subst ...Subst,
) T {
	var result T
	ext := filepath.Ext(fixture)
	if ext == "" {
		fileName := fixture
		if len(subst) > 0 {
			fileName += ".template"
		}
		fileName = FixturePath(t, fileName+".json")
		_, err := os.Stat(fileName)
		if !os.IsNotExist(err) {
			assert.NoError(t, err)
			fixture = fileName
		} else {
			t.Error(fmt.Errorf("fixture %s not found", fixture))
		}
	}
	switch filepath.Ext(fixture) {
	case ".json":
		payload := GetFixture(t, fixture, subst...)
		err := json.Unmarshal(payload, &result)
		assert.NoError(t, err)
	default:
		t.Error(fmt.Errorf("fixture %s not found", fixture))
	}
	return result
}
