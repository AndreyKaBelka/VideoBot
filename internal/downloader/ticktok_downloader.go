package downloader

import (
	"VideoBot/internal/link"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tebeka/selenium"
)

func (d *Service) ExtractCdnUrlFromTiktok(ctx context.Context, link link.Link) (CdnUrl, error) {
	if err := d.wd.Get(link.Link()); err != nil {
		return "", err
	}
	d.logger.Info("Я туууууут")

	err := d.wd.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		source, err := driver.FindElement(selenium.ByXPATH, "//source[@data-index=\"2\"]")
		if err != nil {
			return false, nil
		}
		return source != nil, nil
	}, 5*time.Second)
	d.logger.Info("Я туууууут2")

	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	// Ищем соурс
	source, err := d.wd.FindElement(selenium.ByXPATH, "//source[@data-index=\"2\"]")
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	videoURL, err := source.GetAttribute("src")
	slog.Info("Получил ссылку на скачивание", "link", videoURL)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	return CdnUrl(videoURL), nil
}
