package storage

import (
	"errors"
	"time"
)

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventAlreadyExists = errors.New("event already exists")
)

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	EndTime      time.Time
	Description  *string
	OwnerId      string
	NotifyBefore time.Duration
}
