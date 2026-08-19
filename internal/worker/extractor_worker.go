package worker

import (
	"VideoBot/internal/downloader"
	link2 "VideoBot/internal/link"
	"VideoBot/internal/taskqueue"
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
)

type JobProducer interface {
	SendCdnUrl(ctx context.Context, job taskqueue.CdnUrlArgs) error
}

type Extractor interface {
	ExtractCdnUrlFromInsta(ctx context.Context, link link2.Link) (downloader.CdnUrl, error)
	ExtractCdnUrlFromTiktok(ctx context.Context, link link2.Link) (downloader.CdnUrl, error)
}

type DownloadWorker struct {
	river.WorkerDefaults[taskqueue.LinkJobArgs]
	logger    *slog.Logger
	extractor Extractor
	producer  JobProducer
}

func NewDownloadWorker(logger *slog.Logger, extractor Extractor, producer JobProducer) *DownloadWorker {
	return &DownloadWorker{logger: logger, extractor: extractor, producer: producer}
}

func (w *DownloadWorker) Work(ctx context.Context, job *river.Job[taskqueue.LinkJobArgs]) error {
	link := link2.NewLinkFromArgs(job.Args)

	var videoUrl downloader.CdnUrl
	var err error
	if link.LinkType() == link2.INSTA {
		videoUrl, err = w.extractor.ExtractCdnUrlFromInsta(ctx, link)
	} else if link.LinkType() == link2.TIKTOK {
		w.logger.Info("Начал выгрузку из тикитока")
		videoUrl, err = w.extractor.ExtractCdnUrlFromTiktok(ctx, link)
	}

	if err != nil {
		return fmt.Errorf("не смог взять ссылку: %w", err)
	}

	err = w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
		CdnUrl: videoUrl.String(),
		ChatID: link.ChatId(),
	})
	if err != nil {
		return fmt.Errorf("не смог отправить ссылку: %w", err)
	}
	return nil
}
