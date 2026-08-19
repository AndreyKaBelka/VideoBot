package platform

import (
	"VideoBot/internal/downloader"
	"VideoBot/internal/worker"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func NewRiver(dbpool *pgxpool.Pool, opts *river.Config) (*river.Client[pgx.Tx], error) {
	client, err := river.NewClient(riverpgxv5.New(dbpool), opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func RegisterDownloadWorkers(logger *slog.Logger, downloader *downloader.Service, producer worker.JobProducer) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, worker.NewDownloadWorker(logger, downloader, producer))

	return workers
}

func RegisterCdnUrlWorkers(logger *slog.Logger, sender worker.Sender) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, worker.NewSenderWorker(logger, sender))

	return workers
}
