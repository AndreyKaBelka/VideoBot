package worker

import (
	"VideoBot/internal/downloader"
	"VideoBot/internal/link"
	"VideoBot/internal/taskqueue"
	"context"
	"log/slog"

	"github.com/riverqueue/river"
)

type JobProducer interface {
	SendCdnUrl(ctx context.Context, job taskqueue.CdnUrlArgs) error
}

type Extractor interface {
	ExtractCdnUrlFromInsta(ctx context.Context, lnk link.Link) (downloader.CdnUrl, error)
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
	lnk, err := linkFromArgs(job.Args)
	if err != nil {
		w.logger.Error("Не смог собрать ссылку из джобы", "err", err)
		w.reportError(ctx, job.Args.ChatID, link.ErrNotSupported.Error())
		return nil
	}

	if lnk.LinkType() != link.INSTA {
		w.logger.Error("Неподдерживаемый тип ссылки", "type", lnk.LinkType().Int())
		w.reportError(ctx, lnk.ChatId(), link.ErrNotSupported.Error())
		return nil
	}

	videoUrl, err := w.extractor.ExtractCdnUrlFromInsta(ctx, lnk)
	if err != nil {
		w.logger.Error("Чтото случилось", "err", err)
		w.reportError(ctx, lnk.ChatId(), err.Error())
		return nil
	}

	if err := w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
		CdnUrl: videoUrl.String(),
		ChatID: lnk.ChatId(),
	}); err != nil {
		w.logger.Error("Чтото случилось", "err", err)
		return nil
	}

	return nil
}

func (w *DownloadWorker) reportError(ctx context.Context, chatID int64, msg string) {
	if err := w.producer.SendCdnUrl(ctx, taskqueue.CdnUrlArgs{
		ChatID: chatID,
		Error:  msg,
	}); err != nil {
		w.logger.Error("Чтото случилось", "err", err)
	}
}

// linkFromArgs переводит DTO очереди в доменную ссылку — маппинг живёт здесь,
// чтобы домен ничего не знал о формате джобы.
func linkFromArgs(args taskqueue.LinkJobArgs) (link.Link, error) {
	id, err := link.ParseID(args.ID)
	if err != nil {
		return link.Link{}, err
	}

	linkType, err := link.TypeFromInt(args.LinkType)
	if err != nil {
		return link.Link{}, err
	}

	return link.Restore(id, args.URL, linkType, args.ChatID), nil
}
