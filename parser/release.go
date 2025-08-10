package parser

import (
	"fmt"
	"webhook/structs"
)

type release struct {
	Action  string          `json:"action"`
	Release structs.Release `json:"release"`
	Sender  structs.User    `json:"sender"`
}

func (e *release) handle() (*structs.Webhook, error) {
	if e.Action != "published" {
		return nil, nil
	}

	return &structs.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) published release [%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Release.TagName,
			e.Release.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarURL,
	}, nil
}
