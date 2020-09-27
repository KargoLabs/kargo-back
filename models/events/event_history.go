package models

import "time"

type EventHistory struct {
	Event   Event     `json:"event"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}
