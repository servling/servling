package model

import (
	"time"

	"github.com/servling/servling/ent"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TokenVersion int       `json:"tokenVersion"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreateUserInput struct {
	Username       string
	HashedPassword string
}

func UserFromEnt(user *ent.User) *User {
	return &User{
		ID:           user.ID,
		Name:         user.Name,
		TokenVersion: user.TokenVersion,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
