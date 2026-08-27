package platform

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
)

func NewDriver(seleniumURL string) (selenium.WebDriver, error) {
	ffCaps := firefox.Capabilities{
		Args: []string{"-headless"},
		Prefs: map[string]interface{}{
			"permissions.default.image": 2,
		},
	}

	caps := selenium.Capabilities{
		"browserName": "firefox",
	}
	caps.AddFirefox(ffCaps)

	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось запустить WebDriver: %v", err)
	}

	return wd, nil
}
