package utils

import "github.com/google/uuid"

func UUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
