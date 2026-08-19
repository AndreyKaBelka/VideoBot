package telegram

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func LoggingMiddleware(log *slog.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			start := time.Now()
			if update.Message != nil {
				log.Info("message", "user_id", update.Message.From.ID, "text", update.Message.Text, "sent_at", start.Format("2006-01-02 15:04:05"))
			}

			next(ctx, b, update)
		}
	}
}

func RecoverMiddleware(log *slog.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic recovered", "panic", r, "update_id", update.ID)
				}
			}()

			next(ctx, b, update)
		}
	}
}
