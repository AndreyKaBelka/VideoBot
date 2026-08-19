package link

import (
	"VideoBot/internal/taskqueue"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Save(ctx context.Context, tx pgx.Tx, link Link) error
}

type ProducerService interface {
	SendTaskTx(ctx context.Context, tx pgx.Tx, job taskqueue.LinkJobArgs) error
}

// Beginner открывает транзакцию. Реализуется, например, *pgxpool.Pool.
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct {
	historyRepo Repository
	logger      *slog.Logger
	db          Beginner
	producer    ProducerService
}

func NewService(historyRepo Repository, log *slog.Logger, db Beginner, producer ProducerService) *Service {
	return &Service{
		historyRepo: historyRepo,
		logger:      log,
		db:          db,
		producer:    producer,
	}
}

func (s *Service) Proceed(ctx context.Context, link Link) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.historyRepo.Save(ctx, tx, link); err != nil {
		s.logger.Warn("Не смог сохранить запись в бд", "link", link, "err", err)
		return err
	}

	linkArgs := taskqueue.LinkJobArgs{
		ID:       link.ID().String(),
		URL:      link.Link(),
		LinkType: link.LinkType().Int(),
		ChatID:   link.ChatId(),
	}

	if err := s.producer.SendTaskTx(ctx, tx, linkArgs); err != nil {
		s.logger.Error("Не смог отправить в очередь", "link", link, "err", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("закоммитить транзакцию: %w", err)
	}

	return nil
}
