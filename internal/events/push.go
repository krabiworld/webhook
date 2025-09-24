package events

import (
	"fmt"
	"regexp"
	"strings"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

var linkRe = regexp.MustCompile(`\[([^]]+)]\((https?://[^)]+)\)`)

type Push struct {
	Commits []struct {
		Id      string `json:"id"`
		Url     string `json:"url"`
		Message string `json:"message"`
	} `json:"commits"`
	Created    bool              `json:"created"`
	Deleted    bool              `json:"deleted"`
	Forced     bool              `json:"forced"`
	Ref        string            `json:"ref"`
	Sender     github.User       `json:"sender"`
	Pusher     github.User       `json:"pusher"`
	Repository github.Repository `json:"repository"`
}

func (e *Push) Handle() (*discord.Webhook, error) {
	if !strings.HasPrefix(e.Ref, "refs/heads/") {
		return nil, nil
	}

	action := ""
	if e.Created {
		if e.Forced {
			action = " force"
		}
		action += " created branch"
	} else if e.Deleted {
		action = " deleted branch"
	} else if e.Forced {
		action = " force pushed"
	}

	branch := strings.TrimPrefix(e.Ref, "refs/heads/")

	footer := fmt.Sprintf(
		"\n- [%s](<%s>)%s on [%s](<%s>)/[%s](<%[5]s/tree/%s>)",
		e.Pusher.Name,
		e.Sender.HtmlUrl,
		action,
		e.Repository.Name,
		e.Repository.HtmlUrl,
		branch,
	)

	if len(e.Commits) == 0 {
		return &discord.Webhook{
			Content:   strings.TrimPrefix(footer, "\n"),
			Username:  e.Pusher.Name,
			AvatarUrl: e.Sender.AvatarUrl,
		}, nil
	}

	var commits strings.Builder
	for _, c := range e.Commits {
		lines := strings.Split(c.Message, "\n")
		first := ""
		if len(lines) > 0 {
			first = lines[0]
		}

		commitMsg := first
		if len(lines) > 1 {
			commitMsg += "..."
		}

		cleanMsg := linkRe.ReplaceAllString(commitMsg, "[$1](<$2>)")

		commits.WriteString(fmt.Sprintf("[`%s`](<%s>) %s\n", c.Id[:7], c.Url, cleanMsg))
	}

	limit := 2000 - (len([]rune(footer)) + len("...") + 1)
	if len([]rune(commits.String())) > limit {
		truncated := string([]rune(commits.String())[:limit]) + "..."
		if !strings.HasSuffix(truncated, ">)") {
			lines := strings.Split(truncated, "\n")
			if len(lines) > 1 {
				truncated = strings.Join(lines[:len(lines)-1], "\n")
			}
		}
		commits.Reset()
		commits.WriteString(truncated + "\n")
	}

	return &discord.Webhook{
		Content:   commits.String() + footer,
		Username:  e.Pusher.Name,
		AvatarUrl: e.Sender.AvatarUrl,
	}, nil
}
