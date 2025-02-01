package events

type Repository struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}

type Sender struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

type Author struct {
	Date     string  `json:"date"`
	Email    *string `json:"email"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
}

type Commit struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Url     string `json:"url"`
}
