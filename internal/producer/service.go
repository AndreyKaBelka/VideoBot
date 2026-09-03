package producer

import (
	"VideoBot/internal/link"
	"VideoBot/internal/storage/pg"
	"VideoBot/internal/taskqueue"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Service struct {
	river  *river.Client[pgx.Tx]
	logger *slog.Logger
}

func New(logger *slog.Logger, river *river.Client[pgx.Tx]) *Service {
	return &Service{
		river:  river,
		logger: logger,
	}
}

// SendTask ставит ссылку в очередь в транзакции, открытой вызывающим кодом,
// чтобы запись в БД и джоба коммитились вместе.
func (s *Service) SendTask(ctx context.Context, lnk link.Link) error {
	tx, err := pg.TxFromContext(ctx)
	if err != nil {
		return err
	}

	job := taskqueue.LinkJobArgs{
		ID:       lnk.ID().String(),
		URL:      lnk.Link(),
		LinkType: lnk.LinkType().Int(),
		ChatID:   lnk.ChatId(),
	}

	_, err = s.river.InsertTx(ctx, tx, job, nil)
	return err
}

func (s *Service) SendCdnUrl(ctx context.Context, job taskqueue.CdnUrlArgs) error {
	_, err := s.river.Insert(ctx, job, nil)
	return err
}
