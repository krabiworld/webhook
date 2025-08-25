package handlers

import (
	"net/url"
	"slices"
	"webhook/client"
	"webhook/config"
	"webhook/structs/discord"
	"webhook/structs/github"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type Event interface {
	handle(queries url.Values) (*discord.Webhook, error)
}

type metaEvent struct {
	Repository github.Repository `json:"repository"`
}

func parseEvent[T Event](body []byte, queries url.Values) (*discord.Webhook, error) {
	var e T
	if err := sonic.Unmarshal(body, &e); err != nil {
		return nil, err
	}
	return e.handle(queries)
}

var eventParsers = map[string]func([]byte, url.Values) (*discord.Webhook, error){
	"push":         parseEvent[*push],
	"workflow_run": parseEvent[*workflowRun],
	"star":         parseEvent[*star],
	"fork":         parseEvent[*fork],
	"release":      parseEvent[*release],
}

func Parse(event string, body []byte, queries url.Values, creds discord.Credentials) {
	if len(config.Get().DisabledEvents) > 0 && slices.Contains(config.Get().DisabledEvents, event) {
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

	eventResult, err := parser(body, queries)
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
