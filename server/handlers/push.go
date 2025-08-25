package handlers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"webhook/structs/discord"
	"webhook/structs/github"
)

type push struct {
	Commits    []github.Commit   `json:"commits"`
	Ref        string            `json:"ref"`
	Sender     github.User       `json:"sender"`
	Pusher     github.User       `json:"pusher"`
	Repository github.Repository `json:"repository"`
}

func (e *push) handle(url.Values) (*discord.Webhook, error) {
	if len(e.Commits) == 0 {
		return nil, nil
	}

	linkRe, err := regexp.Compile(`\[([^]]+)]\((https?://[^)]+)\)`)
	if err != nil {
		return nil, err
	}
	mdRe, err := regexp.Compile(`(?m)^\s*#{1,3}\s+`)
	if err != nil {
		return nil, err
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

		cleanLinkMsg := linkRe.ReplaceAllString(commitMsg, "[$1](<$2>)")
		cleanMdMsg := mdRe.ReplaceAllString(cleanLinkMsg, "")

		commits.WriteString(fmt.Sprintf("[`%s`](<%s>) %s\n", c.Id[:7], c.Url, cleanMdMsg))
	}

	branch := strings.TrimPrefix(e.Ref, "refs/heads/")
	footer := fmt.Sprintf(
		"\n- [%s](<%s>) on [%s](<%s>)/[%s](<%s/tree/%s>)",
		e.Pusher.Name,
		e.Sender.HtmlUrl,
		e.Repository.Name,
		e.Repository.HtmlUrl,
		branch,
		e.Repository.HtmlUrl,
		branch,
	)

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
