package events

import (
	"fmt"
	"slices"
	"webhook/internal/config"
	"webhook/internal/context"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"

	"github.com/rs/zerolog/log"
)

type checkRun struct {
	Action     string            `json:"action"`
	CheckRun   github.CheckRun   `json:"check_run"`
	Repository github.Repository `json:"repository"`
}

func (e *checkRun) handle(ctx *context.Context) (*discord.Webhook, error) {
	if e.Action != "completed" || e.CheckRun.Conclusion == "" || e.CheckRun.App.Name == "GitHub Actions" {
		return nil, nil
	}

	if slices.Contains(config.Get().IgnoredChecks, e.CheckRun.App.Name) || slices.Contains(ctx.IgnoredChecks(), e.CheckRun.App.Name) {
		log.Debug().Str("check", e.CheckRun.App.Name).Msg("Ignored check")
		return nil, nil
	}

	emoji := config.Get().SuccessEmoji
	if e.CheckRun.Conclusion == "failure" {
		emoji = config.Get().FailureEmoji
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"%s Check [%s](<%s>) on [%s](<%s>)/[%[6]s](<%[5]s/tree/%[6]s>)",
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
