package github

type Repository struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
	Owner   User   `json:"owner"`
	Private bool   `json:"private"`
}
