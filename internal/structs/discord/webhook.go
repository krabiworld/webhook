package discord

type Webhook struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}
