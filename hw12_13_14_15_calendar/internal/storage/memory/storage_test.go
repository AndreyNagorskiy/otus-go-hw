package memorystorage

import (
	"context"
	"errors"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"testing"
	"time"
)

func TestStorage_CreateEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	event := storage.Event{ID: "1", Title: "Test Event"}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Errorf("CreateEvent() error = %v, want nil", err)
	}

	err = s.CreateEvent(ctx, event)
	if !errors.Is(err, storage.ErrEventAlreadyExists) {
		t.Errorf("CreateEvent() error = %v, want %v", err, storage.ErrEventAlreadyExists)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.CreateEvent(cancelCtx, event)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("CreateEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_GetEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	event := storage.Event{ID: "1", Title: "Test Event"}

	_, err := s.GetEvent(ctx, "1")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("GetEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	_ = s.CreateEvent(ctx, event)
	gotEvent, err := s.GetEvent(ctx, "1")
	if err != nil {
		t.Errorf("GetEvent() error = %v, want nil", err)
	}
	if gotEvent.ID != event.ID || gotEvent.Title != event.Title {
		t.Errorf("GetEvent() = %v, want %v", gotEvent, event)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = s.GetEvent(cancelCtx, "1")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("GetEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	event := storage.Event{ID: "1", Title: "Test Event"}
	updatedEvent := storage.Event{ID: "1", Title: "Updated Event"}

	err := s.UpdateEvent(ctx, updatedEvent)
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("UpdateEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	_ = s.CreateEvent(ctx, event)
	err = s.UpdateEvent(ctx, updatedEvent)
	if err != nil {
		t.Errorf("UpdateEvent() error = %v, want nil", err)
	}
	gotEvent, _ := s.GetEvent(ctx, "1")
	if gotEvent.Title != updatedEvent.Title {
		t.Errorf("After UpdateEvent() title = %v, want %v", gotEvent.Title, updatedEvent.Title)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.UpdateEvent(cancelCtx, updatedEvent)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("UpdateEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_DeleteEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	event := storage.Event{ID: "1", Title: "Test Event"}

	err := s.DeleteEvent(ctx, "1")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("DeleteEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	_ = s.CreateEvent(ctx, event)
	err = s.DeleteEvent(ctx, "1")
	if err != nil {
		t.Errorf("DeleteEvent() error = %v, want nil", err)
	}
	_, err = s.GetEvent(ctx, "1")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("After DeleteEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.DeleteEvent(cancelCtx, "1")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("DeleteEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_GetAllEvents(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	events := []storage.Event{
		{ID: "1", Title: "Event 1"},
		{ID: "2", Title: "Event 2"},
	}

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents() error = %v, want nil", err)
	}
	if len(allEvents) != 0 {
		t.Errorf("GetAllEvents() length = %v, want 0", len(allEvents))
	}

	for _, e := range events {
		_ = s.CreateEvent(ctx, e)
	}
	allEvents, err = s.GetAllEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents() error = %v, want nil", err)
	}
	if len(allEvents) != len(events) {
		t.Errorf("GetAllEvents() length = %v, want %v", len(allEvents), len(events))
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = s.GetAllEvents(cancelCtx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("GetAllEvents() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()
	var wg sync.WaitGroup
	count := 100

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			event := storage.Event{ID: string(rune(i)), Title: "Event"}
			_ = s.CreateEvent(ctx, event)
			_, _ = s.GetEvent(ctx, event.ID)
			_ = s.UpdateEvent(ctx, storage.Event{ID: event.ID, Title: "Updated Event"})
			_ = s.DeleteEvent(ctx, event.ID)
		}(i)
	}
	wg.Wait()

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents() error = %v, want nil", err)
	}
	if len(allEvents) > 0 {
		t.Errorf("After concurrent operations, storage should be empty, got %d events", len(allEvents))
	}
}

func TestStorage_ContextTimeout(t *testing.T) {
	s := NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	time.Sleep(time.Millisecond)

	event := storage.Event{ID: "1", Title: "Test Event"}

	tests := []struct {
		name string
		fn   func() error
		want error
	}{
		{
			name: "CreateEvent",
			fn:   func() error { return s.CreateEvent(ctx, event) },
			want: context.DeadlineExceeded,
		},
		{
			name: "GetEvent",
			fn:   func() error { _, err := s.GetEvent(ctx, "1"); return err },
			want: context.DeadlineExceeded,
		},
		{
			name: "UpdateEvent",
			fn:   func() error { return s.UpdateEvent(ctx, event) },
			want: context.DeadlineExceeded,
		},
		{
			name: "DeleteEvent",
			fn:   func() error { return s.DeleteEvent(ctx, "1") },
			want: context.DeadlineExceeded,
		},
		{
			name: "GetAllEvents",
			fn:   func() error { _, err := s.GetAllEvents(ctx); return err },
			want: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if !errors.Is(err, tt.want) {
				t.Errorf("%s() error = %v, want %v", tt.name, err, tt.want)
			}
		})
	}
}
