package models

const BootstrapToken = "00000000-0000-7000-8000-000000000001"

type Token struct {
	ID    string `gorm:"primaryKey"`
	Token string `gorm:"unique;not null"`
}
