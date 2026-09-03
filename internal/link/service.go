package link

import (
	"context"
	"log/slog"
)

type Repository interface {
	Save(ctx context.Context, link Link) error
}

type ProducerService interface {
	SendTask(ctx context.Context, link Link) error
}

// Transactor выполняет fn в одной транзакции. Как именно транзакция доезжает
// до репозитория и продюсера — деталь реализации, приложению она не видна.
type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	historyRepo Repository
	logger      *slog.Logger
	tx          Transactor
	producer    ProducerService
}

func NewService(historyRepo Repository, log *slog.Logger, tx Transactor, producer ProducerService) *Service {
	return &Service{
		historyRepo: historyRepo,
		logger:      log,
		tx:          tx,
		producer:    producer,
	}
}

func (s *Service) Proceed(ctx context.Context, link Link) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.historyRepo.Save(ctx, link); err != nil {
			s.logger.Warn("Не смог сохранить запись в бд", "link", link.Link(), "err", err)
			return err
		}

		if err := s.producer.SendTask(ctx, link); err != nil {
			s.logger.Error("Не смог отправить в очередь", "link", link.Link(), "err", err)
			return err
		}

		return nil
	})
}
