package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/cybernetlab/swimming-search/config"
	"github.com/cybernetlab/swimming-search/internal/adapter/booking"
	"github.com/cybernetlab/swimming-search/internal/adapter/redis"
	"github.com/cybernetlab/swimming-search/internal/controller/bot"
	"github.com/cybernetlab/swimming-search/internal/controller/bot/commands"
	"github.com/cybernetlab/swimming-search/internal/controller/http"
	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
	"github.com/cybernetlab/swimming-search/pkg/httpserver"
	redislib "github.com/cybernetlab/swimming-search/pkg/redis"
	"github.com/cybernetlab/swimming-search/pkg/router"
)

func Run(ctx context.Context, c config.Config) error {
	// Redis
	redisClient, err := redislib.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redislib.New: %w", err)
	}

	// Booking
	booking := booking.New(c.Booking)

	// Telegram bot
	tg, err := bot.New(c.Bot, commands.Handler)
	if err != nil {
		return fmt.Errorf("bot.New: %w", err)
	}

	// UseCase
	uc := usecase.New(
		redis.New(redisClient),
		booking,
		tg,
	)

	// Metrics
	// httpMetrics := metrics.NewHTTPServer()

	// Start controllers
	// HTTP
	r := router.New()
	http.Router(r)
	httpServer := httpserver.New(r, c.HTTP)

	// Telegram bot
	tg.Run(ctx, uc)

	// Seed store and restart persisted search jobs
	seed(ctx, uc, c.App)
	restartSearches(ctx, uc)

	log.Info().Msg("App started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info().Msg("App got signal to stop")

	// Controllers close
	httpServer.Close()

	// Adapters close
	redisClient.Close(ctx)

	log.Info().Msg("App stopped")

	return nil
}

func seed(ctx context.Context, uc *usecase.UseCase, c config.App) {
	if c.AdminUser == "" {
		return
	}
	_, err := uc.GetUser(ctx, dto.GetUserInput{UserName: c.AdminUser})
	if errors.Is(err, domain.ErrNotFound) {
		_, err := uc.CreateUser(ctx, dto.CreateUserInput{UserName: c.AdminUser, IsAdmin: true})
		if err != nil {
			log.Warn().Msgf("Can't create admin user %s", c.AdminUser)
		}
		log.Info().Msgf("New admin user %s created", c.AdminUser)
		return
	}
	if err != nil {
		log.Error().Err(err).Msgf("Error while creaing admin user %s", c.AdminUser)
	}
}

func restartSearches(ctx context.Context, uc *usecase.UseCase) {
	searches, err := uc.GetSearches(ctx)
	if err == nil {
		if l := len(searches); l > 0 {
			log.Info().Msgf("Restarting %d searches", l)
		}
		for _, search := range searches {
			err = uc.StartSearch(bot.WithChatID(ctx, search.ChatID), dto.StartSearchInput{Search: &search})
			if err != nil {
				log.Warn().Msgf("Can't restart search job: %s", err.Error())
			}
		}
	}
}
