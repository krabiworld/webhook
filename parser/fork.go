package parser

import (
	"fmt"
	"webhook/structs"
)

type fork struct {
	Sender structs.User   `json:"sender"`
	Forkee structs.Forkee `json:"forkee"`
}

func (e *fork) handle() (*structs.Webhook, error) {
	return &structs.Webhook{
		Content: fmt.Sprintf(
			"[%s](<%s>) forked [%s](<%s>)",
			e.Sender.Login,
			e.Sender.HtmlUrl,
			e.Forkee.Name,
			e.Forkee.HtmlUrl,
		),
		Username:  e.Sender.Login,
		AvatarUrl: e.Sender.AvatarURL,
	}, nil
}
