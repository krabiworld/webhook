package github

type CheckRun struct {
	Conclusion string     `json:"conclusion"`
	HtmlUrl    string     `json:"html_url"`
	App        App        `json:"app"`
	CheckSuite CheckSuite `json:"check_suite"`
}
