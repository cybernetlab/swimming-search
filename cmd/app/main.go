package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/cybernetlab/course-progress/config"
	"github.com/cybernetlab/course-progress/internal/app"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/pkg/logger"
	"github.com/cybernetlab/course-progress/pkg/otel"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := domain.WithNodeID(context.Background(), domain.NodeID(c.App.NodeID))

	err = otel.Init(ctx, c.OTEL)
	if err != nil {
		log.Fatal().Err(err).Msg("otel.Init")
	}
	defer otel.Close()

	err = app.Run(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
