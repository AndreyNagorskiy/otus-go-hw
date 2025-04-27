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
	ID           string         `db:"id"`
	Title        string         `db:"title"`
	StartTime    time.Time      `db:"start_time"`
	EndTime      time.Time      `db:"end_time"`
	Description  *string        `db:"description"`
	OwnerID      string         `db:"owner_id"`
	NotifyBefore *time.Duration `db:"notify_before"`
}

type CreateOrUpdateEventParams struct {
	Title        string
	StartTime    time.Time
	EndTime      time.Time
	Description  *string
	OwnerID      string
	NotifyBefore *time.Duration
}
