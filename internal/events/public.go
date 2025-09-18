package events

import (
	"fmt"
	"webhook/internal/config"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

type Public struct {
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

func (e *Public) Handle() (*discord.Webhook, error) {
	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) made [%s](<%s>) public %s",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
			config.Get().HappyEmoji,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
