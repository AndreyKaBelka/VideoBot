package main

import (
	"VideoBot/internal/downloader"
	"VideoBot/internal/platform"
	"VideoBot/internal/producer"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/tebeka/selenium"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load("configs/.env"); err != nil {
		panic("Failed to load .env file")
	}

	// --- конфиг ---
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

	caps := selenium.Capabilities{
		"browserName": "firefox",
	}

	seleniumURL := os.Getenv("SELENIUM_URL")

	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return fmt.Errorf("не удалось запустить WebDriver: %v", err)
	}

	defer wd.Quit()

	downloaderService := downloader.NewDownloader(wd, logger)

	// отдельный клиент без Queues/Workers — только для вставки джоб из воркера,
	// чтобы не завязывать конструирование Workers на клиент, который сам их принимает
	producerRiverClient, err := platform.NewRiver(pool, &river.Config{})
	if err != nil {
		return fmt.Errorf("create river producer client: %w", err)
	}
	producerService := producer.New(logger, producerRiverClient)

	workers := platform.RegisterDownloadWorkers(logger, downloaderService, producerService)

	riverClient, err := platform.NewRiver(pool, &river.Config{
		Queues: map[string]river.QueueConfig{
			"links": {MaxWorkers: 1},
		},
		Workers: workers,
	})
	if err != nil {
		return fmt.Errorf("create river: %w", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		return fmt.Errorf("start river: %w", err)
	}

	<-riverClient.Stopped()
	return nil
}
