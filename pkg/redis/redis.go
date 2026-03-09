package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Addr           string `envconfig:"REDIS_ADDR" required:"true"`
	Password       string `envconfig:"REDIS_PASSWORD"`
	DB             int    `envconfig:"REDIS_DB" default:"0"`
	PersistOnClose bool   `envconfig:"REDIS_PERSIST_ON_CLOSE" default:"false"`
}

type Client struct {
	*redis.Client
	persistOnClose bool
}

func New(c Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	return &Client{Client: client, persistOnClose: c.PersistOnClose}, nil
}

func (c *Client) Close(ctx context.Context) {
	if c.persistOnClose {
		log.Info().Msg("redis: saving database")
		err := c.Client.Save(ctx).Err()
		if err != nil {
			log.Warn().Err(err).Msg("redis: save on close")
		}
	}

	err := c.Client.Close()
	if err != nil {
		log.Error().Err(err).Msg("redis: close")
	}

	log.Info().Msg("redis: closed")
}
