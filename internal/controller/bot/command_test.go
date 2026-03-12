package bot_test

import (
	"context"
	"testing"

	"github.com/cybernetlab/swimming-search/internal/controller/bot"
	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/test/test_tg"
	"github.com/stretchr/testify/require"
)

func New(t *testing.T, ctx context.Context, name string, args string) (bot.Command, *bot.Bot, test_tg.Server) {
	s := test_tg.New(t)
	context.AfterFunc(ctx, func() { s.Server.Close() })
	tg, err := bot.New(bot.Config{URL: s.Server.URL + "/bot%s/%s"}, func(*usecase.UseCase, *bot.Command) {})
	require.NoError(t, err)
	return bot.NewCommand(ctx, tg, name, args), tg, s
}

func TestCommandContext(t *testing.T) {
	cmd, _, _ := New(t, t.Context(), "cmd", "args")
	require.Equal(t, t.Context(), cmd.Context())
}

func TestCommandUser(t *testing.T) {
	u := domain.User{Name: "user"}
	ctx := bot.WithUser(t.Context(), &u)
	cmd, _, _ := New(t, ctx, "cmd", "args")

	user, err := cmd.User()
	require.NoError(t, err)
	require.Equal(t, &u, user)
}

func TestCommandChatID(t *testing.T) {
	ctx := bot.WithChatID(t.Context(), 10)
	cmd, _, _ := New(t, ctx, "cmd", "args")

	chatID, err := cmd.ChatID()
	require.NoError(t, err)
	require.Equal(t, int64(10), chatID)
}

func TestCommandSend(t *testing.T) {
	ctx := bot.WithChatID(t.Context(), 10)
	cmd, _, server := New(t, ctx, "cmd", "args")

	err := cmd.Send("test message")
	require.NoError(t, err)

	msg := <-server.Messages
	require.Equal(t, "10", msg.ChatID)
	require.Equal(t, "test message", msg.Text)
}

func TestCommandSend_WrongContext(t *testing.T) {
	var ctxErr domain.ErrInvalidContext

	cmd, _, _ := New(t, t.Context(), "cmd", "args")
	err := cmd.Send("test message")
	require.ErrorAs(t, err, &ctxErr)
	require.Equal(t, "chatID", ctxErr.Field)
}
