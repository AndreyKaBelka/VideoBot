package worker

import (
	"VideoBot/internal/taskqueue"
	"context"
	"log/slog"

	"github.com/riverqueue/river"
)

type Sender interface {
	SendVideo(ctx context.Context, chatID int64, videoUrl string)
	SendMessage(ctx context.Context, chatID int64, err string)
}

type SenderWorker struct {
	river.WorkerDefaults[taskqueue.CdnUrlArgs]
	logger *slog.Logger
	sender Sender
}

func NewSenderWorker(logger *slog.Logger, sender Sender) *SenderWorker {
	return &SenderWorker{logger: logger, sender: sender}
}

func (w *SenderWorker) Work(ctx context.Context, job *river.Job[taskqueue.CdnUrlArgs]) error {
	if job.Args.Error != "" {
		w.sender.SendMessage(ctx, job.Args.ChatID, job.Args.Error)
		return nil
	}
	w.sender.SendVideo(ctx, job.Args.ChatID, job.Args.CdnUrl)
	return nil
}
