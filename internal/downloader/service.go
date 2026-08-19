package downloader

import (
	"log/slog"

	"github.com/tebeka/selenium"
)

type Service struct {
	wd     selenium.WebDriver
	logger *slog.Logger
}

func NewDownloader(wd selenium.WebDriver, logger *slog.Logger) *Service {
	return &Service{
		wd:     wd,
		logger: logger,
	}
}
