package events

import (
	"encoding/json"
	"log/slog"
	"webhook/internal/client"
	"webhook/internal/structs/discord"
)

type Event interface {
	handle() (*discord.Webhook, error)
}

func parseEvent[T Event](body []byte) (*discord.Webhook, error) {
	var e T
	if err := json.Unmarshal(body, &e); err != nil {
		return nil, err
	}
	return e.handle()
}

var eventParsers = map[string]func([]byte) (*discord.Webhook, error){
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

func Parse(event string, body []byte, creds discord.Credentials) {
	parser, ok := eventParsers[event]
	if !ok {
		slog.Debug("Unknown event", "event", event)
		return
	}

	eventResult, err := parser(body)
	if err != nil {
		slog.Error(err.Error(), "event", event)
		return
	}

	if eventResult == nil {
		return
	}

	err = client.ExecuteWebhook(eventResult, creds)
	if err != nil {
		slog.Error(err.Error(), "event", event)
	}
}
