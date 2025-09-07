package models

type Webhook struct {
	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Token string `gorm:"unique;not null"`
}
