package events

import (
	"fmt"
	"gohook/internal/structs/discord"
	"gohook/internal/structs/github"
)

type Fork struct {
	Forkee struct {
		Name    string `json:"name"`
		HtmlUrl string `json:"html_url"`
	} `json:"forkee"`
	Sender github.User `json:"sender"`
}

func (e *Fork) Handle() (*discord.Webhook, error) {
	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) forked [%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Forkee.Name,
			e.Forkee.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
