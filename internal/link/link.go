package link

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type Type int

const (
	INSTA Type = iota
)

func (t Type) Int() int {
	return int(t)
}

// TypeFromInt восстанавливает тип из числового представления (БД, очередь).
func TypeFromInt(v int) (Type, error) {
	switch Type(v) {
	case INSTA:
		return INSTA, nil
	default:
		return 0, fmt.Errorf("%w: %d", ErrUnknownType, v)
	}
}

type ID uuid.UUID

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// ParseID восстанавливает идентификатор из строкового представления.
func ParseID(s string) (ID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return ID{}, fmt.Errorf("разобрать id ссылки: %w", err)
	}

	return ID(parsed), nil
}

type Link struct {
	id       ID
	link     string
	linkType Type
	chatId   int64
}

func (l *Link) ID() ID {
	return l.id
}

func (l *Link) Link() string {
	return l.link
}

func (l *Link) LinkType() Type {
	return l.linkType
}

func (l *Link) ChatId() int64 {
	return l.chatId
}

var ErrNotSupported = errors.New("not supported link")
var ErrEmptyLink = errors.New("empty link")
var ErrTooLongLink = errors.New("too long link")
var ErrUnknownType = errors.New("unknown link type")

const maxLinkLength = 200

// NewLink создаёт новую ссылку из текста сообщения пользователя.
func NewLink(rawURL string, chatID int64) (Link, error) {
	url := strings.TrimSpace(rawURL)

	if url == "" {
		return Link{}, ErrEmptyLink
	}

	if len(url) > maxLinkLength {
		return Link{}, ErrTooLongLink
	}

	if !isInstagramURL(url) {
		return Link{}, ErrNotSupported
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Link{}, fmt.Errorf("сгенерировать id ссылки: %w", err)
	}

	return Link{
		id:       ID(id),
		link:     url,
		linkType: INSTA,
		chatId:   chatID,
	}, nil
}

// Restore собирает ссылку из уже сохранённого состояния — из очереди или БД,
// минуя проверки, которые она прошла при создании.
func Restore(id ID, rawURL string, linkType Type, chatID int64) Link {
	return Link{
		id:       id,
		link:     rawURL,
		linkType: linkType,
		chatId:   chatID,
	}
}

func isInstagramURL(url string) bool {
	// Регулярное выражение для проверки ссылок Instagram
	pattern := `^(https?:\/\/)?(www\.)?(instagram\.com|instagr\.am)\/([a-zA-Z0-9_\.]+)`

	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}

	return matched
}
