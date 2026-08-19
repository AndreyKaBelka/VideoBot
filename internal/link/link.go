package link

import (
	"VideoBot/internal/taskqueue"
	"errors"
	"regexp"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type Type int

const (
	TIKTOK Type = iota
	INSTA
)

func (t Type) Int() int {
	return int(t)
}

type ID uuid.UUID

func (id ID) String() string {
	return uuid.UUID(id).String()
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

const maxLinkLength = 200

func NewLink(message *models.Update) (Link, error) {
	link := strings.TrimSpace(message.Message.Text)
	chatId := message.Message.Chat.ID

	if link == "" {
		return Link{}, ErrEmptyLink
	}

	if len(link) > maxLinkLength {
		return Link{}, ErrTooLongLink
	}

	return getLinkWithType(link, chatId)
}

func NewLinkFromArgs(args taskqueue.LinkJobArgs) Link {
	return Link{
		id:       ID(uuid.MustParse(args.ID)),
		link:     args.URL,
		linkType: Type(args.LinkType),
		chatId:   args.ChatID,
	}
}

func getLinkWithType(link string, chatId int64) (Link, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Link{}, err
	}
	if isInstagramURL(link) {
		return Link{
			id:       ID(id),
			link:     link,
			linkType: INSTA,
			chatId:   chatId,
		}, nil
	} else if isTikTokURL(link) {
		return Link{
			id:       ID(id),
			link:     link,
			linkType: TIKTOK,
			chatId:   chatId,
		}, nil
	}
	return Link{}, ErrNotSupported
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

func isTikTokURL(url string) bool {
	pattern := `^(https?:\/\/)?(www\.)?(tiktok\.com|vm\.tiktok\.com|vt\.tiktok\.com)\/([@a-zA-Z0-9_\-\.]+\/)?([a-zA-Z0-9_\-\.\/?&=]+)?$`

	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}

	return matched
}
