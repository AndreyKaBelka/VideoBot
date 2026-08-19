package taskqueue

import (
	"github.com/riverqueue/river"
)

type LinkJobArgs struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	LinkType int    `json:"link_type"`
	ChatID   int64  `json:"chat_id"`
}

func (LinkJobArgs) Kind() string { return "link_job" }

func (LinkJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: "links",
	}
}

type CdnUrlArgs struct {
	CdnUrl string `json:"cdn_url"`
	ChatID int64  `json:"chat_id"`
}

func (CdnUrlArgs) Kind() string { return "cdn_url_job" }

func (CdnUrlArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: "cdn_url",
	}
}
