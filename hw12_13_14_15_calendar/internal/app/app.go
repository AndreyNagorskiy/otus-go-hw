package app

import (
	"context"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  logger.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, params storage.CreateOrUpdateEventParams) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
}

type Application interface {
	CreateEvent(ctx context.Context, param storage.CreateOrUpdateEventParams) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger logger.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, param storage.CreateOrUpdateEventParams) error {
	return a.storage.CreateEvent(ctx, param)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetAllEvents(ctx)
}
