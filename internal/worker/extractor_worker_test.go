package worker

import (
	"VideoBot/internal/link"
	"VideoBot/internal/taskqueue"
	"errors"
	"testing"
)

func TestLinkFromArgs(t *testing.T) {
	args := taskqueue.LinkJobArgs{
		ID:       "0199ba39-1c6f-7000-8000-000000000001",
		URL:      "https://instagram.com/reel/Cxyz123/",
		LinkType: link.INSTA.Int(),
		ChatID:   -1001234567890,
	}

	lnk, err := linkFromArgs(args)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if lnk.ID().String() != args.ID {
		t.Errorf("ID() = %s, ожидалось %s", lnk.ID(), args.ID)
	}
	if lnk.Link() != args.URL {
		t.Errorf("Link() = %q, ожидалось %q", lnk.Link(), args.URL)
	}
	if lnk.LinkType() != link.INSTA {
		t.Errorf("LinkType() = %d, ожидалось INSTA", lnk.LinkType().Int())
	}
	if lnk.ChatId() != args.ChatID {
		t.Errorf("ChatId() = %d, ожидалось %d", lnk.ChatId(), args.ChatID)
	}
}

func TestLinkFromArgsRejectsBadInput(t *testing.T) {
	if _, err := linkFromArgs(taskqueue.LinkJobArgs{ID: "не-uuid", LinkType: link.INSTA.Int()}); err == nil {
		t.Error("ожидалась ошибка на некорректном id")
	}

	_, err := linkFromArgs(taskqueue.LinkJobArgs{ID: "0199ba39-1c6f-7000-8000-000000000001", LinkType: 42})
	if !errors.Is(err, link.ErrUnknownType) {
		t.Errorf("ожидалась ErrUnknownType, получено %v", err)
	}
}
