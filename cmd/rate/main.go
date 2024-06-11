package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/irbgeo/rate/internal/coinbase"
	"github.com/irbgeo/rate/internal/controller"
	"github.com/irbgeo/rate/internal/rater"
	"github.com/irbgeo/rate/internal/storage"
	"github.com/irbgeo/rate/pkg/api/http"
	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	Port string `envconfig:"PORT" default:"8080"`

	DBHost string `envconfig:"DB_HOST" default:"localhost"`
	DBPort string `envconfig:"DB_PORT" default:"5432"`
	DBUser string `envconfig:"DB_USER" default:"rate"`
	DBPass string `envconfig:"DB_PASS" default:"rate"`
	DBName string `envconfig:"DB_NAME" default:"rate"`

	RateUpdateInterval time.Duration `envconfig:"RATE_UPDATE_INTERVAL" default:"1m"`
	PairUpdateInterval time.Duration `envconfig:"PAIR_UPDATE_INTERVAL" default:"1m"`
}

func main() {
	var cfg configuration

	err := envconfig.Process("", &cfg)
	if err != nil {
		slog.Error("read config", slog.Any("error", err))
		os.Exit(1)
	}

	stor := storage.NewStorage()

	storStartOpts := storage.StartOpts{
		Username: cfg.DBUser,
		Password: cfg.DBPass,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBName,
	}

	err = stor.Open(storStartOpts)
	if err != nil {
		slog.Error("open storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer stor.Close()

	cb := coinbase.NewService()

	r := rater.NewService(stor, cb)

	rateStartOpts := rater.StartOpts{
		RateUpdateInterval: cfg.RateUpdateInterval,
		PairUpdateInterval: cfg.PairUpdateInterval,
	}

	err = r.Start(rateStartOpts)
	if err != nil {
		slog.Error("start rater", slog.Any("error", err))
		os.Exit(1)
	}

	ctrl := controller.New(r)

	srv := http.NewServer(ctrl)

	go func() {
		err := srv.ListenAndServe(cfg.Port)
		if err != nil {
			slog.Error("listen and serve", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	slog.Info("server started", slog.String("port", cfg.Port))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c

	slog.Info("goodbye!")
}
