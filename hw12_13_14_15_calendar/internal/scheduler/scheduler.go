package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/amqp"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	reminderInterval time.Duration
	cleanupInterval  time.Duration
	l                logger.Logger
	storage          app.Storage
	publisher        amqp.IPublisher
	stopChan         chan struct{}
}

func NewScheduler(
	reminderInterval, cleanupInterval time.Duration,
	l logger.Logger,
	storage app.Storage,
	publisher amqp.IPublisher,
) *Scheduler {
	return &Scheduler{
		reminderInterval: reminderInterval,
		cleanupInterval:  cleanupInterval,
		l:                l,
		storage:          storage,
		publisher:        publisher,
		stopChan:         make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	reminderTicker := time.NewTicker(s.reminderInterval)
	cleanupTicker := time.NewTicker(s.cleanupInterval)
	defer func() {
		reminderTicker.Stop()
		cleanupTicker.Stop()
	}()

	for {
		select {
		case <-reminderTicker.C:
			s.scanAndProcessEvents(ctx)
		case <-cleanupTicker.C:
			s.cleanupOldEvents(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) scanAndProcessEvents(ctx context.Context) {
	s.l.Debug("Scanning for events to notify")

	events, err := s.storage.GetEventsToNotify(ctx)
	if err != nil {
		s.l.Error("Failed to get pending events", slog.String("error", err.Error()))
		return
	}

	for _, event := range events {
		if err := s.processEvent(event); err != nil {
			s.l.Error("Failed to process event", slog.String("event_id", event.ID), slog.String("error", err.Error()))
		}
	}
}

func (s *Scheduler) processEvent(event storage.Event) error {
	s.l.Debug("Processing event", slog.String("event_id", event.ID))

	n := models.Notification{
		ID:        event.ID,
		Title:     event.Title,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		OwnerID:   event.OwnerID,
	}

	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err := s.publisher.Publish("", body, amqp.PublishOptions{Persistent: true}); err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	return nil
}

func (s *Scheduler) cleanupOldEvents(ctx context.Context) {
	s.l.Debug("Cleaning up old events")
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	deletedCount, err := s.storage.DeleteEventsOlderThan(ctx, oneYearAgo)
	if err != nil {
		s.l.Error("Failed to delete old events", slog.String("error", err.Error()))
		return
	}

	s.l.Info(fmt.Sprintf("Successfully deleted %d old events", deletedCount))
}
