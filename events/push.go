package events

type Push struct {
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
	Pusher     Author     `json:"pusher"`
	Sender     Sender     `json:"sender"`
	Commits    []Commit   `json:"commits"`
}
