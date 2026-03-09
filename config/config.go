package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/cybernetlab/course-progress/internal/adapter/booking"
	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/pkg/httpserver"
	"github.com/cybernetlab/course-progress/pkg/logger"
	"github.com/cybernetlab/course-progress/pkg/otel"
	"github.com/cybernetlab/course-progress/pkg/redis"
)

type App struct {
	Name      string `envconfig:"APP_NAME"    required:"true"`
	Version   string `envconfig:"APP_VERSION"`
	NodeID    string `envconfig:"NODE_ID"     required:"true"`
	AdminUser string `envconfig:"ADMIN_USER"`
}

type Config struct {
	App     App
	Logger  logger.Config
	Redis   redis.Config
	HTTP    httpserver.Config
	Bot     bot.Config
	Booking booking.Config
	OTEL    otel.Config
}

func New(defaultVersion string) (Config, error) {
	var config Config

	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		err := godotenv.Load(".env")
		if err != nil {
			return config, fmt.Errorf("godotenv.Load: %w", err)
		}
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	if config.App.Version == "" {
		config.App.Version = defaultVersion
	}

	return config, nil
}
