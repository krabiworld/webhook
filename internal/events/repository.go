package events

import (
	"fmt"
	"slices"
	"strings"
	"webhook/internal/config"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

var supportedRepositoryActions = []string{
	"archived",
	"privatized",
	"renamed",
	"unarchived",
}

type Repository struct {
	Action  string `json:"action"`
	Changes *struct {
		Repository struct {
			Name struct {
				From string `json:"from"`
			} `json:"name"`
		} `json:"repository"`
	}
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

func (e *Repository) Handle() (*discord.Webhook, error) {
	if !slices.Contains(supportedRepositoryActions, e.Action) {
		return nil, nil
	}

	content := strings.Builder{}
	content.WriteString(fmt.Sprintf("[%s](<%s>) ", e.Sender.Login, e.Sender.HtmlUrl))

	switch e.Action {
	case "archived", "unarchived":
		var emoji string
		if e.Action == "archived" {
			emoji = "😭"
		} else if e.Action == "unarchived" {
			emoji = config.Get().HappyEmoji
		}

		content.WriteString(fmt.Sprintf(
			"%s [%s](<%s>) %s",
			e.Action,
			e.Repository.Name,
			e.Repository.HtmlUrl,
			emoji,
		))
	case "privatized":
		content.WriteString(fmt.Sprintf(
			"made [%s](<%s>) private 🖕",
			e.Repository.Name,
			e.Repository.HtmlUrl,
		))
	case "renamed":
		content.WriteString(fmt.Sprintf(
			"renamed [%s](<%s>) to [%s](<%[2]s>)",
			e.Changes.Repository.Name.From,
			e.Repository.HtmlUrl,
			e.Repository.Name,
		))
	}

	return &discord.Webhook{
		Content:   content.String(),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
