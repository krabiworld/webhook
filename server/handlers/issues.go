package handlers

import (
	"fmt"
	"net/url"
	"webhook/structs/discord"
	"webhook/structs/github"
)

type issues struct {
	Action     string            `json:"action"`
	Issue      github.Issue      `json:"issue"`
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

func (e *issues) handle(url.Values) (*discord.Webhook, error) {
	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) %s issue [%s](<%s>) in [%s](<%s>)/[%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Action,
			e.Issue.Title,
			e.Issue.HtmlUrl,
			e.Repository.Owner.Login,
			e.Repository.Owner.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
