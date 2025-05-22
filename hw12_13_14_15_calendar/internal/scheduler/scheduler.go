package scheduler

import (
	"context"
	"fmt"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Scheduler struct {
	reminderInterval time.Duration
	cleanupInterval  time.Duration
	l                logger.Logger
	storage          app.Storage
	stopChan         chan struct{}
}

func NewScheduler(reminderInterval, cleanupInterval time.Duration, l logger.Logger, storage app.Storage) *Scheduler {
	return &Scheduler{
		reminderInterval: reminderInterval,
		cleanupInterval:  cleanupInterval,
		l:                l,
		storage:          storage,
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
	events, err := s.storage.GetEventsToNotify(ctx)
	if err != nil {
		s.l.Error("Failed to get pending events: %v", err)
		return
	}

	for _, event := range events {
		if err := s.processEvent(event); err != nil {
			s.l.Error("Failed to process event %v: %v", event.ID, err)
		}
	}
}

func (s *Scheduler) processEvent(event storage.Event) error {
	//TODO RabbitMQ logic
	s.l.Debug("Processing event: %v", event)
	return nil
}

func (s *Scheduler) cleanupOldEvents(ctx context.Context) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	deletedCount, err := s.storage.DeleteEventsOlderThan(ctx, oneYearAgo)
	if err != nil {
		s.l.Error("Failed to delete old events: %v", err)
		return
	}

	s.l.Info(fmt.Sprintf("Successfully deleted %d old events", deletedCount))
}
