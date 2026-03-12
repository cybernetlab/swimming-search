package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/cybernetlab/swimming-search/internal/domain"
	redislib "github.com/cybernetlab/swimming-search/pkg/redis"
)

type Redis struct {
	redis *redislib.Client
}

var ErrInvalidContext = errors.New("No NodeID in Context")

func New(client *redislib.Client) *Redis {
	return &Redis{redis: client}
}

func (r *Redis) GetUser(ctx context.Context, name string) (domain.User, error) {
	user, err := get[domain.User](ctx, r.redis.Client, "users:"+name)
	if err != nil {
		return user, fmt.Errorf("get[domain.User]: %w", err)
	}
	return user, nil
}

func (r *Redis) GetUsers(ctx context.Context) ([]domain.User, error) {
	keys, err := keys(ctx, r.redis.Client, "users:*")
	if err != nil {
		return []domain.User{}, fmt.Errorf("keys: %w", err)
	}
	users, err := mget[domain.User](ctx, r.redis.Client, keys)
	if err != nil {
		return []domain.User{}, fmt.Errorf("mget: %w", err)
	}
	return users, nil
}

func (r *Redis) PutUser(ctx context.Context, user domain.User) error {
	err := set(ctx, r.redis.Client, "users:"+user.Name, user)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return nil
}

func (r *Redis) DeleteUser(ctx context.Context, name string) error {
	key := fmt.Sprintf("users:%s", name)
	err := delete(ctx, r.redis.Client, key)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *Redis) GetCentres(ctx context.Context) ([]domain.Centre, error) {
	centres, err := get[[]domain.Centre](ctx, r.redis.Client, "centres")
	if err != nil {
		return centres, fmt.Errorf("get[[]domain.Centre]: %w", err)
	}
	return centres, nil
}

func (r *Redis) PutCentres(ctx context.Context, centres []domain.Centre) error {
	err := set(ctx, r.redis.Client, "centres", centres)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return nil
}

func (r *Redis) GetSearch(ctx context.Context, userName string) (domain.Search, domain.NodeID, error) {
	var search domain.Search
	var node domain.NodeID

	keys, err := keys(ctx, r.redis.Client, fmt.Sprintf("searches:%s:*", userName))
	if err != nil {
		return search, node, fmt.Errorf("keys: %w", err)
	}
	if len(keys) != 1 {
		return search, node, domain.ErrNotFound
	}
	parts := strings.SplitN(keys[0], ":", 3)
	node = domain.NodeID(parts[2])
	search, err = get[domain.Search](ctx, r.redis.Client, keys[0])
	if err != nil {
		return search, node, fmt.Errorf("get[domain.Search]: %w", err)
	}
	return search, node, nil
}

func (r *Redis) GetSearches(ctx context.Context, nodeID domain.NodeID) ([]domain.Search, error) {
	node, err := domain.ContextNodeID(ctx)
	if err != nil {
		return []domain.Search{}, fmt.Errorf("domain.ContextNodeID: %w", err)
	}
	keys, err := keys(ctx, r.redis.Client, fmt.Sprintf("searches:*:%s", node))
	if err != nil {
		return []domain.Search{}, fmt.Errorf("keys: %w", err)
	}
	searches, err := mget[domain.Search](ctx, r.redis.Client, keys)
	if err != nil {
		return []domain.Search{}, fmt.Errorf("mget: %w", err)
	}
	return searches, nil
}

func (r *Redis) PutSearch(ctx context.Context, search domain.Search) error {
	node, err := domain.ContextNodeID(ctx)
	if err != nil {
		return fmt.Errorf("domain.ContextNodeID: %w", err)
	}
	key := fmt.Sprintf("searches:%s:%s", search.UserName, node)
	err = set(ctx, r.redis.Client, key, search)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return nil
}

func (r *Redis) DeleteSearch(ctx context.Context, search domain.Search) error {
	node, err := domain.ContextNodeID(ctx)
	if err != nil {
		return fmt.Errorf("domain.ContextNodeID: %w", err)
	}
	key := fmt.Sprintf("searches:%s:%s", search.UserName, node)
	err = delete(ctx, r.redis.Client, key)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func get[T []domain.Centre | domain.User | domain.Search](ctx context.Context, r *redis.Client, key string) (T, error) {
	var result T
	val, err := r.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, domain.ErrNotFound
		}
		return result, fmt.Errorf("r.Get: %w", err)
	}
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return result, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return result, nil
}

func set[T []domain.Centre | domain.User | domain.Search](ctx context.Context, r *redis.Client, key string, val T) error {
	payload, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	err = r.Set(ctx, key, payload, 0).Err()
	if err != nil {
		return fmt.Errorf("r.Set: %w", err)
	}
	return nil
}

func keys(ctx context.Context, r *redis.Client, pattern string) ([]string, error) {
	result := []string{}

	var cursor uint64
	for {
		keys, cursor, err := r.Scan(ctx, cursor, pattern, 30).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return []string{}, nil
			}
			return []string{}, fmt.Errorf("r.Scan: %w", err)
		}
		result = append(result, keys...)
		if cursor == 0 {
			return result, nil
		}
	}
}

func mget[T domain.Search | domain.User](ctx context.Context, r *redis.Client, keys []string) ([]T, error) {
	var result []T
	var entry T

	val, err := r.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, domain.ErrNotFound
		}
		return result, fmt.Errorf("r.MGet: %w", err)
	}
	for _, v := range val {
		err = json.Unmarshal([]byte(v.(string)), &entry)
		if err != nil {
			return result, fmt.Errorf("json.Unmarshal: %w", err)
		}
		result = append(result, entry)
	}
	return result, nil
}

func delete(ctx context.Context, rdb *redis.Client, key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("r.Del: %w", err)
	}
	return nil
}
