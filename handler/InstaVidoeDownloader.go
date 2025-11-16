package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const outputName = "test.mp4"

type InstagramDownloader struct {
	wd selenium.WebDriver
}

func NewInstagramDownloader() (*InstagramDownloader, error) {
	// Настройка Selenium сервера
	caps := selenium.Capabilities{
		"browserName": "firefox",
	}

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		return nil, fmt.Errorf("не удалось запустить WebDriver: %v", err)
	}

	return &InstagramDownloader{wd: wd}, nil
}

func (id *InstagramDownloader) Close() {
	err := id.wd.Quit()
	if err != nil {
		slog.Error("Не смог закрыть селениум", "error", err)
	}
}

func (id *InstagramDownloader) StartDownloadReel(reelURL, filename string) (io.ReadCloser, error) {
	// Переходим на страницу рила
	if err := id.wd.Get(reelURL); err != nil {
		return nil, err
	}

	defer id.Close()

	err := id.wd.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		scripts, _ := driver.FindElements(selenium.ByXPATH, "//script[contains(., 'video_versions')]")
		if scripts != nil {
			return true, nil
		}
		return false, nil
	}, 30*time.Second)

	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	// Ищем видео элемент
	script, err := id.wd.FindElement(selenium.ByXPATH, "//script[contains(., 'video_versions')]")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	jsonObj, err := id.FireScript(script)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	videoURL, err := extractVideoURLFromJSON(jsonObj)
	slog.Info("Получил ссылку на скачивание", "link", videoURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить видео элемент: %v", err)
	}

	return id.downloadVideo(videoURL, filename)
}

func (id *InstagramDownloader) downloadVideo(videoURL, filename string) (io.ReadCloser, error) {
	// Создаем HTTP клиент
	client := &http.Client{}

	// Создаем запрос
	req, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовки как у браузера
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://www.instagram.com/")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	slog.Info("Закончил загрузку видео")

	return resp.Body, nil
}

func (id *InstagramDownloader) FireScript(element selenium.WebElement) (interface{}, error) {
	script := `
						var element = arguments[0];
						return element.textContent || element.innerText || element.innerHTML;
					`
	return id.wd.ExecuteScript(script, []interface{}{element})
}

func extractVideoURLFromJSON(result interface{}) (string, error) {
	jsonStr, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("неверный формат данных")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	// Используем безопасное извлечение с проверками
	videoURL, err := getNestedString(data, "require", "0", "3", "0", "__bbox", "require", "0", "3", "1", "__bbox", "result", "data", "xdt_api__v1__media__shortcode__web_info", "items", "0", "video_versions", "0", "url")

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(videoURL, `\/`, `/`), nil
}

// Вспомогательная функция для безопасного извлечения вложенных данных
func getNestedString(data map[string]interface{}, keys ...string) (string, error) {
	var current interface{} = data

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		case []interface{}:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(v) {
				return "", fmt.Errorf("неверный индекс: %s", key)
			}
			current = v[idx]
		default:
			return "", fmt.Errorf("неверный тип на ключе %s", key)
		}

		if current == nil {
			return "", fmt.Errorf("ключ не найден: %s", key)
		}
	}

	result, ok := current.(string)
	if !ok {
		return "", fmt.Errorf("значение не является строкой")
	}

	return result, nil
}
