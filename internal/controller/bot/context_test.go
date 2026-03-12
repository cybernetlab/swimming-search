package bot_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/controller/bot"
	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestContextChatID(t *testing.T) {
	var ctxErr domain.ErrInvalidContext

	id, err := bot.ContextChatID(t.Context())
	require.ErrorAs(t, err, &ctxErr)
	require.Equal(t, "chatID", ctxErr.Field)

	ctx := bot.WithChatID(t.Context(), 10)
	id, err = bot.ContextChatID(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(10), id)
}

func TestContextUser(t *testing.T) {
	var ctxErr domain.ErrInvalidContext

	id, err := bot.ContextUser(t.Context())
	require.ErrorAs(t, err, &ctxErr)
	require.Equal(t, "user", ctxErr.Field)

	user := domain.User{Name: "test"}
	ctx := bot.WithUser(t.Context(), &user)
	id, err = bot.ContextUser(ctx)
	require.NoError(t, err)
	require.Equal(t, &user, id)
}
