package events

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

var ignoredPullRequestActions = []string{
	"labeled",
	"synchronize",
}

type pullRequest struct {
	Action      string            `json:"action"`
	PullRequest github.Issue      `json:"pull_request"`
	Repository  github.Repository `json:"repository"`
	Sender      github.User       `json:"sender"`
}

func (e *pullRequest) handle() (*discord.Webhook, error) {
	if strings.Contains(e.Action, "_") {
		return nil, nil
	}

	if slices.Contains(ignoredPullRequestActions, e.Action) {
		slog.Debug("Ignoring action", "action", e.Action)
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
