package events

import (
	"fmt"
	"strings"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

type pullRequest struct {
	Action      string             `json:"action"`
	PullRequest github.PullRequest `json:"pull_request"`
	Repository  github.Repository  `json:"repository"`
	Sender      github.User        `json:"sender"`
}

func (e *pullRequest) handle() (*discord.Webhook, error) {
	if strings.Contains(e.Action, "_") || e.Action == "synchronize" {
		return nil, nil
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) %s pull request [%s](<%s>) in [%s](<%s>)/[%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Action,
			e.PullRequest.Title,
			e.PullRequest.HtmlUrl,
			e.Repository.Owner.Login,
			e.Repository.Owner.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
