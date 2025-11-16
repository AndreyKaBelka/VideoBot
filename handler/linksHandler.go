package handler

import (
	"fmt"
	"regexp"
)

type LinkType int

const (
	TIKTOK LinkType = iota
	INSTA
)

type Link struct {
	link     string
	linkType LinkType
}

func IsInstagramURL(url string) bool {
	// Регулярное выражение для проверки ссылок Instagram
	pattern := `^(https?:\/\/)?(www\.)?(instagram\.com|instagr\.am)\/([a-zA-Z0-9_\.]+)`

	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}

	return matched
}

func IsTikTokURL(url string) bool {
	pattern := `^(https?:\/\/)?(www\.)?(tiktok\.com|vm\.tiktok\.com|vt\.tiktok\.com)\/([@a-zA-Z0-9_\-\.]+\/)?([a-zA-Z0-9_\-\.\/?&=]+)?$`

	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}

	return matched
}

func GetLinkType(link string) (Link, error) {
	if IsInstagramURL(link) {
		return Link{
			link:     link,
			linkType: INSTA,
		}, nil
	} else if IsTikTokURL(link) {
		if IsInstagramURL(link) {
			return Link{
				link:     link,
				linkType: TIKTOK,
			}, nil
		}
	}
	return Link{}, &NotSupportedLinkError{
		message: "Not supported link type: %s",
		link:    link,
	}

}

type NotSupportedLinkError struct {
	message string
	link    string
}

func (e *NotSupportedLinkError) Error() string {
	return fmt.Sprintf(e.message, e.link)
}
