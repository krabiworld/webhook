package events

import (
	"fmt"
	"gohook/internal/structs/discord"
	"gohook/internal/structs/github"
)

type Release struct {
	Action  string `json:"action"`
	Release struct {
		HtmlUrl string `json:"html_url"`
		TagName string `json:"tag_name"`
	} `json:"release"`
	Sender github.User `json:"sender"`
}

func (e *Release) Handle() (*discord.Webhook, error) {
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
