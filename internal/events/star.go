package events

import (
	"fmt"
	"time"
	"webhook/internal/config"
	"webhook/internal/debouncer"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

type Star struct {
	Action     string            `json:"action"`
	Sender     github.User       `json:"sender"`
	Repository github.Repository `json:"repository"`
}

func (e *Star) Handle() (*discord.Webhook, error) {
	if e.Action != "created" {
		return nil, nil
	}

	ok := debouncer.Debounce("star", e.Sender.Login, e.Repository.Name, time.Hour*24)
	if !ok {
		return nil, nil
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) starred [%s](<%s>) %s",
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
