package events

import (
	"fmt"
	"time"
	"webhook/config"
	"webhook/context"
	"webhook/debouncer"
	"webhook/structs/discord"
	"webhook/structs/github"

	"github.com/rs/zerolog/log"
)

type star struct {
	Action     string            `json:"action"`
	Sender     github.User       `json:"sender"`
	Repository github.Repository `json:"repository"`
}

func (e *star) handle(*context.Context) (*discord.Webhook, error) {
	if e.Action != "created" {
		return nil, nil
	}

	ok := debouncer.Debounce("star", e.Sender.Login, e.Repository.Name, time.Hour*24)
	if !ok {
		log.Debug().Str("repository", e.Repository.Name).Str("username", e.Sender.Login).Msg("Event is currently debounced")
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
