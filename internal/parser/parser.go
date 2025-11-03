package parser

import (
	"fmt"
	"webhook/internal/client"
	"webhook/internal/events"
	"webhook/internal/structs/discord"

	"encoding/json"
)

type Event interface {
	Handle() (*discord.Webhook, error)
}

func parseEvent[T Event](body []byte) (*discord.Webhook, error) {
	var e T
	if err := json.Unmarshal(body, &e); err != nil {
		return nil, err
	}
	return e.Handle()
}

var eventParsers = map[string]func([]byte) (*discord.Webhook, error){
	"check_run":     parseEvent[*events.CheckRun],
	"fork":          parseEvent[*events.Fork],
	"issue_comment": parseEvent[*events.IssueComment],
	"issues":        parseEvent[*events.Issues],
	"public":        parseEvent[*events.Public],
	"pull_request":  parseEvent[*events.PullRequest],
	"push":          parseEvent[*events.Push],
	"release":       parseEvent[*events.Release],
	"repository":    parseEvent[*events.Repository],
	"star":          parseEvent[*events.Star],
	"workflow_run":  parseEvent[*events.WorkflowRun],
}

func Parse(event string, body []byte, creds discord.Credentials) {
	parser, ok := eventParsers[event]
	if !ok {
		return
	}

	eventResult, err := parser(body)
	if err != nil {
		fmt.Println("Parsing error:", err)
		return
	}

	if eventResult == nil {
		return
	}

	err = client.ExecuteWebhook(eventResult, creds)
	if err != nil {
		fmt.Println("Client error:", err)
	}
}
