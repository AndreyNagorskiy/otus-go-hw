package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

func makeCreateOrUpdateEventParams() storage.CreateOrUpdateEventParams {
	description := "Test Description"
	notifyBefore := time.Duration(10) * time.Minute

	return storage.CreateOrUpdateEventParams{
		Title:        "Test Event",
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(1 * time.Hour),
		Description:  &description,
		OwnerID:      "123e4567-e89b-12d3-a456-426614174000",
		NotifyBefore: &notifyBefore,
	}
}

func TestStorage_CreateEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()

	params := makeCreateOrUpdateEventParams()
	err := s.CreateEvent(ctx, params)
	if err != nil {
		t.Errorf("CreateEvent() error = %v, want nil", err)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.CreateEvent(cancelCtx, params)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("CreateEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_GetEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()

	_, err := s.GetEvent(ctx, "nonexistent-id")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("GetEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	params := makeCreateOrUpdateEventParams()
	err = s.CreateEvent(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(allEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(allEvents))
	}
	eventID := allEvents[0].ID

	gotEvent, err := s.GetEvent(ctx, eventID)
	if err != nil {
		t.Errorf("GetEvent() error = %v, want nil", err)
	}
	if gotEvent.Title != params.Title {
		t.Errorf("GetEvent() = %v, want title %v", gotEvent, params.Title)
	}
	if gotEvent.StartTime != params.StartTime {
		t.Errorf("GetEvent() = %v, want start time %v", gotEvent.StartTime, params.StartTime)
	}
	if gotEvent.EndTime != params.EndTime {
		t.Errorf("GetEvent() = %v, want end time %v", gotEvent.EndTime, params.EndTime)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = s.GetEvent(cancelCtx, eventID)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("GetEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()

	params := makeCreateOrUpdateEventParams()
	err := s.CreateEvent(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	eventID := allEvents[0].ID

	err = s.UpdateEvent(ctx, storage.Event{ID: "nonexistent-id", Title: "Updated"})
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("UpdateEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	updatedEvent := storage.Event{
		ID:        eventID,
		Title:     "Updated Event",
		StartTime: allEvents[0].StartTime,
		EndTime:   allEvents[0].EndTime,
		OwnerID:   allEvents[0].OwnerID,
	}
	err = s.UpdateEvent(ctx, updatedEvent)
	if err != nil {
		t.Errorf("UpdateEvent() error = %v, want nil", err)
	}

	gotEvent, err := s.GetEvent(ctx, eventID)
	if err != nil {
		t.Fatal(err)
	}
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

	err := s.DeleteEvent(ctx, "nonexistent-id")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("DeleteEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	params := makeCreateOrUpdateEventParams()
	err = s.CreateEvent(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	eventID := allEvents[0].ID

	err = s.DeleteEvent(ctx, eventID)
	if err != nil {
		t.Errorf("DeleteEvent() error = %v, want nil", err)
	}

	_, err = s.GetEvent(ctx, eventID)
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("After DeleteEvent() error = %v, want %v", err, storage.ErrEventNotFound)
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.DeleteEvent(cancelCtx, eventID)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("DeleteEvent() with canceled context error = %v, want %v", err, context.Canceled)
	}
}

func TestStorage_GetAllEvents(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()

	allEvents, err := s.GetAllEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents() error = %v, want nil", err)
	}
	if len(allEvents) != 0 {
		t.Errorf("GetAllEvents() length = %v, want 0", len(allEvents))
	}

	testEvents := []storage.CreateOrUpdateEventParams{
		makeCreateOrUpdateEventParams(),
		{
			Title:     "Event 2",
			StartTime: time.Now().Add(2 * time.Hour),
			EndTime:   time.Now().Add(3 * time.Hour),
			OwnerID:   "owner-2",
		},
	}

	for _, e := range testEvents {
		err := s.CreateEvent(ctx, e)
		if err != nil {
			t.Fatal(err)
		}
	}

	allEvents, err = s.GetAllEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents() error = %v, want nil", err)
	}
	if len(allEvents) != len(testEvents) {
		t.Errorf("GetAllEvents() length = %v, want %v", len(allEvents), len(testEvents))
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
			params := storage.CreateOrUpdateEventParams{
				Title:     fmt.Sprintf("Event %d", i),
				StartTime: time.Now().Add(time.Duration(i) * time.Minute),
				EndTime:   time.Now().Add(time.Duration(i+1) * time.Minute),
				OwnerID:   fmt.Sprintf("owner-%d", i),
			}
			err := s.CreateEvent(ctx, params)
			if err != nil {
				t.Logf("CreateEvent failed: %v", err)
				return
			}

			allEvents, err := s.GetAllEvents(ctx)
			if err != nil {
				t.Logf("GetAllEvents failed: %v", err)
				return
			}

			var eventID string
			for _, e := range allEvents {
				if e.Title == params.Title && e.OwnerID == params.OwnerID {
					eventID = e.ID
					break
				}
			}

			if eventID == "" {
				t.Log("Event not found after creation")
				return
			}

			_, err = s.GetEvent(ctx, eventID)
			if err != nil {
				t.Logf("GetEvent failed: %v", err)
				return
			}

			updatedEvent := storage.Event{
				ID:        eventID,
				Title:     fmt.Sprintf("Updated Event %d", i),
				StartTime: params.StartTime,
				EndTime:   params.EndTime,
				OwnerID:   params.OwnerID,
			}
			err = s.UpdateEvent(ctx, updatedEvent)
			if err != nil {
				t.Logf("UpdateEvent failed: %v", err)
				return
			}

			err = s.DeleteEvent(ctx, eventID)
			if err != nil {
				t.Logf("DeleteEvent failed: %v", err)
				return
			}
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

	params := makeCreateOrUpdateEventParams()

	tests := []struct {
		name string
		fn   func() error
		want error
	}{
		{
			name: "CreateEvent",
			fn:   func() error { return s.CreateEvent(ctx, params) },
			want: context.DeadlineExceeded,
		},
		{
			name: "GetEvent",
			fn:   func() error { _, err := s.GetEvent(ctx, "nonexistent-id"); return err },
			want: context.DeadlineExceeded,
		},
		{
			name: "UpdateEvent",
			fn: func() error {
				return s.UpdateEvent(ctx, storage.Event{
					ID:        "nonexistent-id",
					Title:     "Test Event",
					StartTime: time.Now(),
					EndTime:   time.Now().Add(time.Hour),
					OwnerID:   "test-owner",
				})
			},
			want: context.DeadlineExceeded,
		},
		{
			name: "DeleteEvent",
			fn:   func() error { return s.DeleteEvent(ctx, "nonexistent-id") },
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
