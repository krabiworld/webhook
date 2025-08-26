package handlers

import (
	"slices"
	"webhook/client"
	"webhook/config"
	"webhook/context"
	"webhook/structs/discord"
	"webhook/structs/github"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type Event interface {
	handle(ctx *context.Context) (*discord.Webhook, error)
}

type metaEvent struct {
	Repository github.Repository `json:"repository"`
}

func parseEvent[T Event](body []byte, ctx *context.Context) (*discord.Webhook, error) {
	var e T
	if err := sonic.Unmarshal(body, &e); err != nil {
		return nil, err
	}
	return e.handle(ctx)
}

var eventParsers = map[string]func([]byte, *context.Context) (*discord.Webhook, error){
	"check_run":     parseEvent[*checkRun],
	"fork":          parseEvent[*fork],
	"issue_comment": parseEvent[*issueComment],
	"issues":        parseEvent[*issues],
	"pull_request":  parseEvent[*pullRequest],
	"push":          parseEvent[*push],
	"release":       parseEvent[*release],
	"star":          parseEvent[*star],
	"workflow_run":  parseEvent[*workflowRun],
}

func Parse(event string, body []byte, ctx *context.Context, creds discord.Credentials) {
	if len(config.Get().DisabledEvents) > 0 && slices.Contains(config.Get().DisabledEvents, event) || slices.Contains(ctx.IgnoredEvents(), event) {
		log.Debug().Str("event", event).Msg("Ignoring event")
		return
	}

	parser, ok := eventParsers[event]
	if !ok {
		log.Debug().Str("event", event).Msg("Unknown event")
		return
	}

	var meta metaEvent
	if err := sonic.Unmarshal(body, &meta); err != nil {
		log.Error().Err(err).Msg("Failed to parse meta")
		return
	}

	if config.Get().IgnorePrivateRepositories && meta.Repository.Private {
		log.Debug().Str("repo", meta.Repository.Name).Msg("Ignoring private repository")
		return
	}

	if len(config.Get().IgnoredRepositories) > 0 && slices.Contains(config.Get().IgnoredRepositories, meta.Repository.Name) {
		log.Debug().Str("repo", meta.Repository.Name).Msg("Ignoring repository")
		return
	}

	eventResult, err := parser(body, ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	if eventResult == nil {
		return
	}

	err = client.ExecuteWebhook(eventResult, creds)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
