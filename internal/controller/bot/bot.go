package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token string `envconfig:"BOT_TOKEN" required:"true"`
	Debug bool   `envconfig:"BOT_DEBUG" default:"false"`
	URL   string
}

type Bot struct {
	api     *tgbotapi.BotAPI
	handler CommandHandler
}

type CommandHandler func(*usecase.UseCase, *Command)

func New(c Config, handler CommandHandler) (*Bot, error) {
	var botAPI *tgbotapi.BotAPI
	var err error

	if c.URL != "" {
		botAPI, err = tgbotapi.NewBotAPIWithAPIEndpoint(c.Token, c.URL)
	} else {
		botAPI, err = tgbotapi.NewBotAPI(c.Token)
	}
	if err != nil {
		return &Bot{}, fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}
	botAPI.Debug = c.Debug
	return &Bot{api: botAPI, handler: handler}, nil
}

func (b *Bot) Send(ctx context.Context, text string) error {
	chatID, err := ContextChatID(ctx)
	if err != nil {
		return fmt.Errorf("ContextChatID: %w", err)
	}
	msg := tgbotapi.NewMessage(int64(chatID), text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = b.api.Send(msg)
	return err
}

func (b *Bot) Sendf(ctx context.Context, text string, args ...any) error {
	return b.Send(ctx, fmt.Sprintf(text, args...))
}

func (b *Bot) Reply(ctx context.Context, text string) {
	err := b.Send(ctx, text)
	if err != nil {
		log.Error().Err(err).Msg("b.Send")
	}
}

func (b *Bot) Replyf(ctx context.Context, text string, args ...any) {
	err := b.Sendf(ctx, text, args...)
	if err != nil {
		log.Error().Err(err).Msg("b.Send")
	}
}

func processUpdate(ctx context.Context, bot *Bot, uc *usecase.UseCase, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		ctx = WithChatID(ctx, update.Message.Chat.ID)
		user, err := uc.GetUser(ctx, dto.GetUserInput{UserName: update.Message.From.UserName})
		if err != nil {
			if !errors.Is(err, domain.ErrNotFound) {
				log.Error().Err(err).Msgf("Error getting user %s", update.Message.From.UserName)
			}
			bot.Reply(ctx, "You are not authorized to use this bot")
			return
		}
		ctx = WithUser(ctx, &user.User)

		log.Debug().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		command := Command{
			Name: update.Message.Command(),
			Args: update.Message.CommandArguments(),
			bot:  bot,
			ctx:  ctx,
		}
		bot.handler(uc, &command)
	}
}

func (bot *Bot) Run(ctx context.Context, uc *usecase.UseCase) {
	go func() {
		log.Info().Msgf("Authorized on account %s", bot.api.Self.UserName)

		for {
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60
			updates := bot.api.GetUpdatesChan(u)

		UPDATES:
			for {
				select {
				case <-ctx.Done():
					return
				case update, ok := <-updates:
					if ok {
						processUpdate(ctx, bot, uc, update)
					} else {
						log.Info().Msg("Disconnected from telegramm")
						break UPDATES
					}
				}
			}
		}
	}()
}
