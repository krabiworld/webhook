package events

type WorkflowRun struct {
	Action      string `json:"action"`
	WorkflowRun struct {
		HeadBranch *string `json:"head_branch"`
		Conclusion *string `json:"conclusion"`
		HtmlUrl    string  `json:"html_url"`
	} `json:"workflow_run"`
	Workflow struct {
		Name string `json:"name"`
	} `json:"workflow"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}
