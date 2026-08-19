package telegram

import (
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
)

func New(token string, h *Handler, log *slog.Logger) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(h.HandleDefault),
		bot.WithMiddlewares(
			RecoverMiddleware(log),
			LoggingMiddleware(log),
		),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	h.sender.SetBot(b)
	h.RegisterHandlers(b)

	return b, nil
}
