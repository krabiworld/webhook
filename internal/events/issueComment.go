package events

import (
	"fmt"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

type IssueComment struct {
	Action  string `json:"action"`
	Comment struct {
		Body    string `json:"body"`
		HtmlUrl string `json:"html_url"`
	} `json:"comment"`
	Issue      github.Issue      `json:"issue"`
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

func (e *IssueComment) Handle() (*discord.Webhook, error) {
	return &discord.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) %s comment on issue [%s](<%s>) in [%s](<%s>)/[%s](<%s>)",
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
