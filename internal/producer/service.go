package producer

import (
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

func (s *Service) SendTaskTx(ctx context.Context, tx pgx.Tx, job taskqueue.LinkJobArgs) error {
	_, err := s.river.InsertTx(ctx, tx, job, nil)
	return err
}

func (s *Service) SendCdnUrl(ctx context.Context, job taskqueue.CdnUrlArgs) error {
	_, err := s.river.Insert(ctx, job, nil)
	return err
}
