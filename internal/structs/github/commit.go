package github

type Commit struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	Message string `json:"message"`
}
