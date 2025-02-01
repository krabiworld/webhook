package events

type Fork struct {
	Forkee struct {
		Name    string `json:"name"`
		HtmlUrl string `json:"html_url"`
	} `json:"forkee"`
	Sender Sender `json:"sender"`
}
