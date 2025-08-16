package handlers

import (
	"fmt"
	"slices"
	"webhook/config"
	"webhook/structs/discord"
	"webhook/structs/github"
)

type workflowRun struct {
	Action      string             `json:"action"`
	Workflow    github.Workflow    `json:"workflow"`
	WorkflowRun github.WorkflowRun `json:"workflowRun"`
	Repository  github.Repository  `json:"repository"`
}

func (e *workflowRun) handle() (*discord.Webhook, error) {
	if e.Action != "completed" || e.WorkflowRun.Conclusion == "" {
		return nil, nil
	}

	if slices.Contains(config.Get().IgnoredWorkflows, e.Workflow.Name) {
		return nil, nil
	}

	emoji := config.Get().SuccessEmoji
	if e.WorkflowRun.Conclusion == "failure" {
		emoji = config.Get().FailureEmoji
	}

	return &discord.Webhook{
		Content: fmt.Sprintf(
			"%s Workflow [%s](<%s>) on [%s](<%s>)/[%[6]s](<%[5]s/tree/%[6]s>)",
			emoji,
			e.WorkflowRun.Conclusion,
			e.WorkflowRun.HtmlUrl,
			e.Repository.Name,
			e.Repository.HtmlUrl,
			e.WorkflowRun.HeadBranch,
		),
		Username:  e.Workflow.Name,
		AvatarUrl: e.Repository.Owner.AvatarURL,
	}, nil
}
