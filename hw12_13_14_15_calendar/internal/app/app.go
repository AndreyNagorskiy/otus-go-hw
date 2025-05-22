package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  logger.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, params storage.CreateOrUpdateEventParams) (*storage.Event, error)
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
	GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error)
	DeleteEventsOlderThan(ctx context.Context, cutoffTime time.Time) (int64, error)
	GetEventsToNotify(ctx context.Context) ([]storage.Event, error)
}

type Application interface {
	CreateEvent(ctx context.Context, param storage.CreateOrUpdateEventParams) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
	GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error)
	GetEventsForDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetEventsForWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error)
	GetEventsForMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error)
}

func New(logger logger.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, param storage.CreateOrUpdateEventParams) (*storage.Event, error) {
	event, err := a.storage.CreateEvent(ctx, param)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExists) {
			a.logger.Info("Event already exists", slog.String("error", err.Error()))
		}

		a.logger.Error("Failed to create event", slog.String("error", err.Error()))
	}

	return event, err
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	err := a.storage.UpdateEvent(ctx, event)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			a.logger.Info("Event not found", slog.String("error", err.Error()))
		}

		a.logger.Error("Failed to update event", slog.String("error", err.Error()))
	}

	return err
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			a.logger.Info("Event not found", slog.String("error", err.Error()))
		}

		a.logger.Error("Failed to delete event", slog.String("error", err.Error()))
	}

	return err
}

func (a *App) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	event, err := a.storage.GetEvent(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			a.logger.Info("Event not found", slog.String("error", err.Error()))
		}

		a.logger.Error("Failed to get event", slog.String("error", err.Error()))
	}

	return event, err
}

func (a *App) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.storage.GetAllEvents(ctx)
	if err != nil {
		a.logger.Error("Failed to get all events", slog.String("error", err.Error()))
	}

	return events, err
}

func (a *App) GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	events, err := a.storage.GetEventsByPeriod(ctx, start, end)
	if err != nil {
		a.logger.Error("Failed to get events by period",
			slog.String("start", start.String()),
			slog.String("end", end.String()),
			slog.String("error", err.Error()))
	}

	return events, nil
}

func (a *App) GetEventsForDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.AddDate(0, 0, 1)

	return a.GetEventsByPeriod(ctx, start, end)
}

func (a *App) GetEventsForWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.AddDate(0, 0, -1)
	}
	end := weekStart.AddDate(0, 0, 7)

	return a.GetEventsByPeriod(ctx, weekStart, end)
}

func (a *App) GetEventsForMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	start := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	end := start.AddDate(0, 1, 0)

	return a.GetEventsByPeriod(ctx, start, end)
}
