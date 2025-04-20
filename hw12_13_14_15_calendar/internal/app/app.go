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
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetAllEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger logger.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}
