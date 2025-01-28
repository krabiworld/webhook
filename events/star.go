package events

type Star struct {
	Action       string       `json:"action"`
	StarredAt    interface{}  `json:"starred_at"`
	Repository   Repository   `json:"repository"`
	Organization Organization `json:"organization"`
	Sender       Sender       `json:"sender"`
}
