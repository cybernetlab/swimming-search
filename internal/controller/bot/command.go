package bot

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/cybernetlab/course-progress/internal/domain"
)

type Command struct {
	Name string
	Args string
	bot  *Bot
	ctx  context.Context
}

func (c *Command) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Command) User() (*domain.User, error) {
	return ContextUser(c.ctx)
}

func (c *Command) ChatID() (int64, error) {
	return ContextChatID(c.ctx)
}

func (c *Command) Send(text string) error {
	return c.bot.Send(c.ctx, text)
}

func (c *Command) Sendf(text string, args ...any) error {
	return c.bot.Sendf(c.ctx, text, args...)
}

func (c *Command) Reply(msg string) {
	err := c.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("cmd.Send")
	}
}

func (c *Command) Replyf(msg string, args ...any) {
	err := c.Sendf(msg, args...)
	if err != nil {
		log.Error().Err(err).Msg("cmd.Send")
	}
}

func (c *Command) Error(msg string, err error, errMsg string) {
	log.Error().Err(err).Msg(errMsg)
	err = c.Send(fmt.Sprintf("%s: internal server error", msg))
	if err != nil {
		log.Error().Err(err).Msg("cmd.Send")
	}
}
