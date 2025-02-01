package events

type Star struct {
	Action     string     `json:"action"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}
