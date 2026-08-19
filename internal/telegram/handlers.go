package telegram

import (
	"VideoBot/internal/link"
	"context"
	"errors"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	log    *slog.Logger
	links  LinkService
	sender *Service
}

func NewHandler(log *slog.Logger, links LinkService, sender *Service) *Handler {
	return &Handler{
		log:    log,
		links:  links,
		sender: sender,
	}
}

func (h *Handler) RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.HandleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, h.HandleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/", bot.MatchTypePrefix, h.HandleUnsupportedCommand)
	b.RegisterHandler(bot.HandlerTypeMessageText, "https://", bot.MatchTypePrefix, h.HandleUrl)
}

func (h *Handler) HandleDefault(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	h.sender.SendMessage(ctx, update.Message.Chat.ID, "/help для возможных команд")
}

func (h *Handler) handleLinkError(ctx context.Context, b *bot.Bot, update *models.Update, err error) {
	switch {
	case errors.Is(err, link.ErrTooLongLink):
		h.sender.SendMessage(ctx, update.Message.Chat.ID, "Слишком длинная ссылка")
	case errors.Is(err, link.ErrEmptyLink):
		h.sender.SendMessage(ctx, update.Message.Chat.ID, "Пустая ссылка")
	case errors.Is(err, link.ErrNotSupported):
		h.sender.SendMessage(ctx, update.Message.Chat.ID, "Не смог понять что за ссылка")
	default:
		h.sender.SendMessage(ctx, update.Message.Chat.ID, "Что-то мне плоховато, повторите запрос позже")
		h.log.Error("Ошибка:", "error", err)
	}
}

func (h *Handler) HandleUnsupportedCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.sender.SendMessage(ctx, update.Message.Chat.ID, "Не понимаю эту команду. Напишите /help")
}

func (h *Handler) HandleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.sender.SendMessage(ctx, update.Message.Chat.ID, "Доступные команды:\n/start — начать")
}

func (h *Handler) HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.sender.SendMessage(ctx, update.Message.Chat.ID,
		"Привет! Этот бот позволяет загружать видео из ТикТока и Instagram. Просто отправь ссылку в чат!")
}
