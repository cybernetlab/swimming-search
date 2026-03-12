package bot

import (
	"context"

	"github.com/cybernetlab/swimming-search/internal/domain"
)

type (
	chatIDKeyType int
	userKeyType   int
)

var (
	chatIdKey chatIDKeyType
	userKey   userKeyType
)

func ContextChatID(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(chatIdKey).(int64)
	if !ok {
		return 0, domain.NewErrInvalidContext("chatID")
	}
	return id, nil
}

func WithChatID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, chatIdKey, id)
}

func ContextUser(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(userKey).(*domain.User)
	if !ok {
		return nil, domain.NewErrInvalidContext("user")
	}
	return user, nil
}

func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
