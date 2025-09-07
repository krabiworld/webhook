package dtos

type Webhook struct {
	Name string `json:"name" validate:"required"`
}
