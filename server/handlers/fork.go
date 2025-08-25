package handlers

import (
	"fmt"
	"net/url"
	"webhook/structs/discord"
	"webhook/structs/github"
)

type fork struct {
	Sender github.User   `json:"sender"`
	Forkee github.Forkee `json:"forkee"`
}

func (e *fork) handle(url.Values) (*discord.Webhook, error) {
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
