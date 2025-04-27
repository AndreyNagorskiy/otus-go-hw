package app

import (
	"context"
	"errors"
	"log/slog"

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
}

type Application interface {
	CreateEvent(ctx context.Context, param storage.CreateOrUpdateEventParams) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
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
