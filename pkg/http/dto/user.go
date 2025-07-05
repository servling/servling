package dto

import (
	"time"

	"github.com/servling/servling/pkg/model"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TokenVersion int       `json:"tokenVersion"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func UserFromModel(user *model.User) *User {
	return &User{
		ID:           user.ID,
		Name:         user.Name,
		TokenVersion: user.TokenVersion,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
