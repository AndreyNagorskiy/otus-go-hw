package memorystorage

import (
	"context"
	"sync"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, params storage.CreateOrUpdateEventParams) (*storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		id := uuid.New().String()

		if _, exists := s.events[id]; exists {
			return nil, storage.ErrEventAlreadyExists
		}

		event := storage.Event{
			ID:           id,
			Title:        params.Title,
			StartTime:    params.StartTime,
			EndTime:      params.EndTime,
			Description:  params.Description,
			OwnerID:      params.OwnerID,
			NotifyBefore: params.NotifyBefore,
		}

		s.events[event.ID] = event
		return &event, nil
	}
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()

		event, exists := s.events[id]
		if !exists {
			return nil, storage.ErrEventNotFound
		}
		return &event, nil
	}
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, exists := s.events[event.ID]; !exists {
			return storage.ErrEventNotFound
		}
		s.events[event.ID] = event
		return nil
	}
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, exists := s.events[id]; !exists {
			return storage.ErrEventNotFound
		}
		delete(s.events, id)
		return nil
	}
}

func (s *Storage) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()

		events := make([]storage.Event, 0, len(s.events))
		for _, e := range s.events {
			events = append(events, e)
		}
		return events, nil
	}
}
