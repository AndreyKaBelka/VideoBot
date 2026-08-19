package main

import (
	"VideoBot/internal/link"
	"VideoBot/internal/platform"
	"VideoBot/internal/producer"
	"VideoBot/internal/telegram"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load("configs/.env"); err != nil {
		panic("Failed to load .env file")
	}

	// --- конфиг ---
	token := os.Getenv("TOKEN")
	if token == "" {
		return fmt.Errorf("TOKEN is not set")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	pool, err := platform.NewPG(dsn)
	if err != nil {
		return fmt.Errorf("create pg pool: %w", err)
	}
	defer pool.Close()

	sender := telegram.NewSender(log)

	workers := platform.RegisterCdnUrlWorkers(log, sender)

	riverClient, err := platform.NewRiver(pool, &river.Config{
		Queues: map[string]river.QueueConfig{
			"cdn_url": {MaxWorkers: 1},
		},
		Workers: workers,
	})
	if err != nil {
		return fmt.Errorf("create river: %w", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		return fmt.Errorf("start river: %w", err)
	}
	defer func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		if err := riverClient.Stop(stopCtx); err != nil {
			log.Error("stop river", "err", err)
		}
	}()

	producerRiverClient, err := platform.NewRiver(pool, &river.Config{})
	if err != nil {
		return fmt.Errorf("create river producer client: %w", err)
	}

	linkRepo := link.NewRepo()
	producerService := producer.New(log, producerRiverClient)
	linkService := link.NewService(linkRepo, log, pool, producerService)
	handler := telegram.NewHandler(log, linkService, sender)

	b, err := telegram.New(token, handler, log)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}

	// --- запуск ---
	log.Info("bot starting")
	b.Start(ctx)
	log.Info("bot stopped")

	return nil
}
