package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	OwnerID   string    `json:"ownerId"`
}
