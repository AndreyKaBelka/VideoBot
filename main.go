package main

import (
	"VideoBot/handler"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(loggerMiddleware),
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.Start(ctx)
}

func loggerMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			slog.Info("message", "user_id", update.Message.From.ID, "text", update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		videoCh := make(chan string)
		defer close(videoCh)
		webDriver, err := handler.NewInstagramDownloader()
		if err != nil {
			slog.Error("Не удалось создать вебДрайвер", "error", err)
		}
		defer webDriver.Close()
		link := update.Message.Text

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Начал загрузку видео...",
		})
		if err != nil {
			slog.Error("Не смог отправить сообщение", "error", err)
		}

		go func() {
			data, err := webDriver.StartDownloadReel(link, "test.mp4")
			if err != nil {
				slog.Error("Не удалось скачать видео", "error", err)
			}
			videoCh <- data
		}()
		video := <-videoCh
		_, err = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID: update.Message.Chat.ID,
			Video: &models.InputFileString{
				Data: video,
			},
		})
		if err != nil {
			slog.Error("Ошибка отправки видео", "err", err)
		}
	}
}
