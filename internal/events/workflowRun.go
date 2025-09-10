package events

import (
	"fmt"
	"log/slog"
	"slices"
	"webhook/internal/config"
	"webhook/internal/structs/discord"
	"webhook/internal/structs/github"
)

var ignoredWorkflows = []string{
	"CodeQL",
	"Dependabot Updates",
	"Automatic Dependency Submission",
}

type workflowRun struct {
	Action   string `json:"action"`
	Workflow struct {
		Name string `json:"name"`
	} `json:"workflow"`
	WorkflowRun struct {
		Conclusion string `json:"conclusion"`
		HtmlUrl    string `json:"html_url"`
		HeadBranch string `json:"head_branch"`
	} `json:"workflow_run"`
	Repository github.Repository `json:"repository"`
}

func (e *workflowRun) handle() (*discord.Webhook, error) {
	if e.Action != "completed" || e.WorkflowRun.Conclusion == "" {
		return nil, nil
	}

	if slices.Contains(ignoredWorkflows, e.Workflow.Name) {
		slog.Debug("Ignoring workflow", "workflow", e.Workflow.Name)
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
		AvatarUrl: e.Repository.Owner.AvatarUrl,
	}, nil
}
