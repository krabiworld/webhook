package events

import (
	"fmt"
	"webhook/context"
	"webhook/structs/discord"
	"webhook/structs/github"
)

type release struct {
	Action  string         `json:"action"`
	Release github.Release `json:"release"`
	Sender  github.User    `json:"sender"`
}

func (e *release) handle(*context.Context) (*discord.Webhook, error) {
	if e.Action != "published" {
		return nil, nil
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) published release [%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Release.TagName,
			e.Release.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
