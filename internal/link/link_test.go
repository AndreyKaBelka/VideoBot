package link

import (
	"errors"
	"strings"
	"testing"
)

func TestNewLink(t *testing.T) {
	const chatID int64 = -1001234567890

	tests := []struct {
		name    string
		rawURL  string
		wantErr error
		wantURL string
	}{
		{
			name:    "reel",
			rawURL:  "https://www.instagram.com/reel/Cxyz123/",
			wantURL: "https://www.instagram.com/reel/Cxyz123/",
		},
		{
			name:    "без схемы",
			rawURL:  "instagram.com/p/Cxyz123/",
			wantURL: "instagram.com/p/Cxyz123/",
		},
		{
			name:    "короткий домен",
			rawURL:  "http://instagr.am/user",
			wantURL: "http://instagr.am/user",
		},
		{
			name:    "пробелы обрезаются",
			rawURL:  "  https://instagram.com/reel/Cxyz123/\n",
			wantURL: "https://instagram.com/reel/Cxyz123/",
		},
		{
			name:    "пустая строка",
			rawURL:  "",
			wantErr: ErrEmptyLink,
		},
		{
			name:    "только пробелы",
			rawURL:  "   \t\n",
			wantErr: ErrEmptyLink,
		},
		{
			name:    "слишком длинная",
			rawURL:  "https://instagram.com/reel/" + strings.Repeat("a", maxLinkLength),
			wantErr: ErrTooLongLink,
		},
		{
			name:    "чужой сайт",
			rawURL:  "https://youtube.com/watch?v=xyz",
			wantErr: ErrNotSupported,
		},
		{
			name:    "инстаграм не в начале строки",
			rawURL:  "смотри https://instagram.com/reel/Cxyz123/",
			wantErr: ErrNotSupported,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lnk, err := NewLink(tt.rawURL, chatID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ожидалась ошибка %v, получено %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if lnk.Link() != tt.wantURL {
				t.Errorf("Link() = %q, ожидалось %q", lnk.Link(), tt.wantURL)
			}
			if lnk.LinkType() != INSTA {
				t.Errorf("LinkType() = %d, ожидалось INSTA", lnk.LinkType().Int())
			}
			if lnk.ChatId() != chatID {
				t.Errorf("ChatId() = %d, ожидалось %d", lnk.ChatId(), chatID)
			}
			if lnk.ID() == (ID{}) {
				t.Error("ID() пустой, ожидался сгенерированный uuid")
			}
		})
	}
}

func TestNewLinkGeneratesUniqueIDs(t *testing.T) {
	const url = "https://instagram.com/reel/Cxyz123/"

	first, err := NewLink(url, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	second, err := NewLink(url, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if first.ID() == second.ID() {
		t.Error("две ссылки получили одинаковый ID")
	}
}

func TestTypeFromInt(t *testing.T) {
	got, err := TypeFromInt(INSTA.Int())
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got != INSTA {
		t.Errorf("TypeFromInt(0) = %d, ожидалось INSTA", got.Int())
	}

	if _, err := TypeFromInt(42); !errors.Is(err, ErrUnknownType) {
		t.Errorf("ожидалась ErrUnknownType, получено %v", err)
	}
}

func TestParseIDRoundTrip(t *testing.T) {
	original, err := NewLink("https://instagram.com/reel/Cxyz123/", 7)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	parsed, err := ParseID(original.ID().String())
	if err != nil {
		t.Fatalf("ParseID вернул ошибку: %v", err)
	}
	if parsed != original.ID() {
		t.Errorf("ParseID вернул %s, ожидалось %s", parsed, original.ID())
	}

	if _, err := ParseID("не-uuid"); err == nil {
		t.Error("ожидалась ошибка на некорректном uuid")
	}
}

func TestRestore(t *testing.T) {
	id, err := ParseID("0199ba39-1c6f-7000-8000-000000000001")
	if err != nil {
		t.Fatalf("ParseID вернул ошибку: %v", err)
	}

	lnk := Restore(id, "https://instagram.com/reel/Cxyz123/", INSTA, 42)

	if lnk.ID() != id || lnk.LinkType() != INSTA || lnk.ChatId() != 42 {
		t.Errorf("Restore вернул неожиданную ссылку: %+v", lnk)
	}
}
