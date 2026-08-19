package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service struct {
	bot    *bot.Bot
	logger *slog.Logger
}

func NewSender(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

// SetBot привязывает *bot.Bot к отправителю. Вызывается после bot.New,
// т.к. сам bot.Bot создаётся из Handler, которому Service нужен заранее.
func (h *Service) SetBot(b *bot.Bot) {
	h.bot = b
}

func (h *Service) SendMessage(ctx context.Context, chatID int64, text string) {
	if _, err := h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		h.logger.Error("failed to send message", "err", err, "chat_id", chatID)
	}
}

func (h *Service) SendVideo(ctx context.Context, chatID int64, videoUrl string) {
	if _, err := h.bot.SendVideo(ctx, &bot.SendVideoParams{
		ChatID: chatID,
		Video: &models.InputFileString{
			Data: videoUrl,
		},
	}); err != nil {
		h.logger.Error("failed to send video", "err", err, "chat_id", chatID)
	}
}
