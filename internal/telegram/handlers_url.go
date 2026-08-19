package telegram

import (
	"VideoBot/internal/link"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type LinkService interface {
	Proceed(ctx context.Context, lnk link.Link) error
}

func (h *Handler) HandleUrl(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	lnk, err := link.NewLink(update)
	if err != nil {
		h.handleLinkError(ctx, b, update, err)
		return
	}

	if err := h.links.Proceed(ctx, lnk); err != nil {
		h.sender.SendMessage(ctx, update.Message.Chat.ID, "Что-то пошло не так")
		return
	}

	h.sender.SendMessage(ctx, update.Message.Chat.ID, "Начал загрузку...")
}
