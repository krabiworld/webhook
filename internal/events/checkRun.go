package events

import (
	"fmt"
	"gohook/internal/config"
	"gohook/internal/structs/discord"
	"gohook/internal/structs/github"
	"slices"
)

var ignoredChecks = []string{
	"Dependabot",
	"GitHub Actions",
	"GitHub Advanced Security",
}

type CheckRun struct {
	Action   string `json:"action"`
	CheckRun struct {
		Conclusion string `json:"conclusion"`
		HtmlUrl    string `json:"html_url"`
		App        struct {
			Name string `json:"name"`
		} `json:"app"`
		CheckSuite struct {
			HeadBranch string `json:"head_branch"`
		} `json:"check_suite"`
	} `json:"check_run"`
	Repository github.Repository `json:"repository"`
}

func (e *CheckRun) Handle() (*discord.Webhook, error) {
	if e.Action != "completed" || e.CheckRun.Conclusion == "" {
		return nil, nil
	}

	if slices.Contains(ignoredChecks, e.CheckRun.App.Name) {
		return nil, nil
	}

	emoji := config.Get().SuccessEmoji
	if e.CheckRun.Conclusion == "failure" {
		emoji = config.Get().FailureEmoji
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"%s Check [%s](<%s>) on [%s](<%s>)/[%s](<%[5]s/tree/%s>)",
			emoji,
			e.CheckRun.Conclusion,
			e.CheckRun.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
			e.CheckRun.CheckSuite.HeadBranch,
		),
		Username:  e.CheckRun.App.Name,
		AvatarUrl: e.Repository.Owner.AvatarUrl,
	}, nil
}
