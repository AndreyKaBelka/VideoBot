package downloader

import (
	"VideoBot/internal/link"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type CdnUrl string

func (c CdnUrl) String() string {
	return string(c)
}

func (d *Service) ExtractCdnUrlFromInsta(ctx context.Context, link link.Link) (CdnUrl, error) {
	if err := d.wd.Get(link.Link()); err != nil {
		return "", err
	}

	err := d.wd.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		scripts, _ := driver.FindElements(selenium.ByXPATH, "//script[contains(., 'video_versions')]")
		if scripts != nil {
			return true, nil
		}
		return false, nil
	}, 30*time.Second)

	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	// Ищем видео элемент
	script, err := d.wd.FindElement(selenium.ByXPATH, "//script[contains(., 'video_versions')]")
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	jsonObj, err := d.fireScript(script)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	videoURL, err := extractVideoURLFromJSON(jsonObj)
	slog.Info("Получил ссылку на скачивание", "link", videoURL)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	return CdnUrl(videoURL), nil
}

func (d *Service) fireScript(element selenium.WebElement) (interface{}, error) {
	script := `
						var element = arguments[0];
						return element.textContent || element.innerText || element.innerHTML;
					`
	return d.wd.ExecuteScript(script, []interface{}{element})
}

func extractVideoURLFromJSON(result interface{}) (string, error) {
	jsonStr, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("неверный формат данных")
	}

	re := regexp.MustCompile(`(?m)"video_versions":\[\{[^}]*"url":"(https?:\\/\\/[^"]+)"`)

	videoURL := re.FindStringSubmatch(jsonStr)

	return strings.ReplaceAll(videoURL[1], `\/`, `/`), nil
}
