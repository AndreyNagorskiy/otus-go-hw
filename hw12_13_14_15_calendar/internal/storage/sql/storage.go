package sqlstorage

import (
	"context"
	"fmt"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateEvent(ctx context.Context, params storage.CreateOrUpdateEventParams) (*storage.Event, error) {
	query := `
		INSERT INTO events (title, start_time, end_time, description, owner_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING`

	_, err := s.db.Exec(ctx, query, params.Title, params.StartTime, params.EndTime, params.Description, params.OwnerID,
		params.NotifyBefore)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	return nil, nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	query := `
		SELECT id, title, start_time, end_time, description, owner_id, notify_before
		FROM events 
		WHERE id = $1`

	var event storage.Event

	err := s.db.QueryRow(ctx, query, id).Scan(&event.ID, &event.Title, &event.StartTime, &event.EndTime,
		&event.Description, &event.OwnerID, &event.NotifyBefore)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `
		UPDATE events
		SET title = $2,
		start_time = $3,
		end_time = $4,
		description = $5,
		owner_id = $6,
		notify_before = $7
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, event.ID, event.Title, event.StartTime, event.EndTime, event.Description,
		event.OwnerID, event.NotifyBefore)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `
		DELETE FROM events
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	query := `
		SELECT id, title, start_time, end_time, description, owner_id, notify_before
		FROM events`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.StartTime, &event.EndTime, &event.Description,
			&event.OwnerID, &event.NotifyBefore); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return events, nil
}
