package github

type User struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}
