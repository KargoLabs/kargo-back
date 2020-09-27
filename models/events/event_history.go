package models

import (
	users "kargo-back/models/users"
	"time"
)

type EventHistory struct {
	Event    Event          `json:"event"`
	Message  string         `json:"message"`
	UserType users.UserType `json:"user_type"`
	Date     time.Time      `json:"date"`
}
