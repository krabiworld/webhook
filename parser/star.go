package parser

import (
	"fmt"
	"time"
	"webhook/config"
	"webhook/jail"
	"webhook/structs"

	"github.com/rs/zerolog/log"
)

type star struct {
	Action     string             `json:"action"`
	Sender     structs.User       `json:"sender"`
	Repository structs.Repository `json:"repository"`
}

func (e *star) handle() (*structs.Webhook, error) {
	if e.Action != "created" {
		return nil, nil
	}

	ok := jail.Trap("star", e.Sender.Login, e.Repository.Name, time.Hour*24)
	if !ok {
		log.Debug().Str("repository", e.Repository.Name).Str("username", e.Sender.Login).Msg("User in jail")
		return nil, nil
	}

	return &structs.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) starred [%s](<%s>) %s",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
			config.Get().HappyEmoji,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarURL,
	}, nil
}
