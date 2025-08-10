package structs

type WorkflowRun struct {
	Conclusion string `json:"conclusion"`
	HtmlUrl    string `json:"html_url"`
	HeadBranch string `json:"head_branch"`
}
