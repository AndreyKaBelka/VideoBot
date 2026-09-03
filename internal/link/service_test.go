package link

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

type stubTransactor struct {
	committed bool
	calls     int
}

func (s *stubTransactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	s.calls++
	if err := fn(ctx); err != nil {
		return err
	}
	s.committed = true
	return nil
}

type stubRepo struct {
	saved []Link
	err   error
}

func (s *stubRepo) Save(_ context.Context, link Link) error {
	if s.err != nil {
		return s.err
	}
	s.saved = append(s.saved, link)
	return nil
}

type stubProducer struct {
	sent []Link
	err  error
}

func (s *stubProducer) SendTask(_ context.Context, link Link) error {
	if s.err != nil {
		return s.err
	}
	s.sent = append(s.sent, link)
	return nil
}

func newTestLink(t *testing.T) Link {
	t.Helper()

	lnk, err := NewLink("https://instagram.com/reel/Cxyz123/", 100)
	if err != nil {
		t.Fatalf("не удалось создать ссылку: %v", err)
	}

	return lnk
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServiceProceedCommits(t *testing.T) {
	repo := &stubRepo{}
	producer := &stubProducer{}
	tx := &stubTransactor{}
	lnk := newTestLink(t)

	if err := NewService(repo, discardLogger(), tx, producer).Proceed(context.Background(), lnk); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if tx.calls != 1 {
		t.Errorf("WithinTx вызван %d раз, ожидался 1", tx.calls)
	}
	if !tx.committed {
		t.Error("транзакция не закоммичена")
	}
	if len(repo.saved) != 1 || repo.saved[0].ID() != lnk.ID() {
		t.Errorf("в репозиторий сохранено %v, ожидалась ссылка %s", repo.saved, lnk.ID())
	}
	if len(producer.sent) != 1 || producer.sent[0].ID() != lnk.ID() {
		t.Errorf("в очередь отправлено %v, ожидалась ссылка %s", producer.sent, lnk.ID())
	}
}

func TestServiceProceedRollsBackOnRepoError(t *testing.T) {
	wantErr := errors.New("бд недоступна")
	repo := &stubRepo{err: wantErr}
	producer := &stubProducer{}
	tx := &stubTransactor{}

	err := NewService(repo, discardLogger(), tx, producer).Proceed(context.Background(), newTestLink(t))

	if !errors.Is(err, wantErr) {
		t.Fatalf("ожидалась ошибка %v, получено %v", wantErr, err)
	}
	if tx.committed {
		t.Error("транзакция закоммичена, хотя запись не удалась")
	}
	if len(producer.sent) != 0 {
		t.Error("джоба отправлена в очередь, хотя запись в бд не удалась")
	}
}

func TestServiceProceedRollsBackOnProducerError(t *testing.T) {
	wantErr := errors.New("очередь недоступна")
	repo := &stubRepo{}
	producer := &stubProducer{err: wantErr}
	tx := &stubTransactor{}

	err := NewService(repo, discardLogger(), tx, producer).Proceed(context.Background(), newTestLink(t))

	if !errors.Is(err, wantErr) {
		t.Fatalf("ожидалась ошибка %v, получено %v", wantErr, err)
	}
	if tx.committed {
		t.Error("транзакция закоммичена, хотя отправка в очередь не удалась")
	}
}
