package worker

import (
	"VideoBot/internal/downloader"
	link2 "VideoBot/internal/link"
	"VideoBot/internal/taskqueue"
	"context"
	"log/slog"

	"github.com/riverqueue/river"
)

type JobProducer interface {
	SendCdnUrl(ctx context.Context, job taskqueue.CdnUrlArgs) error
}

type Extractor interface {
	ExtractCdnUrlFromInsta(ctx context.Context, link link2.Link) (downloader.CdnUrl, error)
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

	if link.LinkType() != link2.INSTA {
		w.logger.Error("Неподдерживаемый тип ссылки", "type", link.LinkType().Int())
		if err := w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
			ChatID: link.ChatId(),
			Error:  link2.ErrNotSupported.Error(),
		}); err != nil {
			w.logger.Error("Чтото случилось", "err", err)
		}
		return nil
	}

	videoUrl, err := w.extractor.ExtractCdnUrlFromInsta(ctx, link)

	if err != nil {
		w.logger.Error("Чтото случилось", "err", err)
		err = w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
			ChatID: link.ChatId(),
			Error:  err.Error(),
		})
		return nil
	}

	err = w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
		CdnUrl: videoUrl.String(),
		ChatID: link.ChatId(),
	})
	if err != nil {
		w.logger.Error("Чтото случилось", "err", err)
		return nil
	}
	return nil
}
