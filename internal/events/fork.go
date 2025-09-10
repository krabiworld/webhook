package events

import (
	"fmt"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

type fork struct {
	Forkee struct {
		Name    string `json:"name"`
		HtmlUrl string `json:"html_url"`
	} `json:"forkee"`
	Sender github.User `json:"sender"`
}

func (e *fork) handle() (*discord.Webhook, error) {
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
