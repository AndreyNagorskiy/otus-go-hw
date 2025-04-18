package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/helpers"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	app       app.Application
	validator *validator.Validate
}

type createOrUpdateEventRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	StartTime    string  `json:"startTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime      string  `json:"endTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
	OwnerID      string  `json:"ownerId" validate:"required,uuid"`
	NotifyBefore *int    `json:"notifyBefore" validate:"omitempty,min=0"`
}

func NewEventHandler(app app.Application) *EventHandler {
	v := helpers.GetValidator()

	return &EventHandler{
		app:       app,
		validator: v,
	}
}

func (e *EventHandler) prepareForCreateOrUpdate(
	w http.ResponseWriter,
	r *http.Request,
) (*storage.CreateOrUpdateEventParams, error) {
	var req createOrUpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, Error("Invalid request payload"))
		return nil, err
	}

	if err := e.validator.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			RespondWithJSON(w, http.StatusBadRequest, ValidationError(validateErr))
			return nil, err
		}
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, Error("Invalid start time format"))
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, Error("Invalid end time format"))
		return nil, err
	}

	var notifyBefore *time.Duration
	if req.NotifyBefore != nil {
		duration := time.Duration(*req.NotifyBefore) * time.Minute
		notifyBefore = &duration
	}

	return &storage.CreateOrUpdateEventParams{
		Title:        req.Title,
		StartTime:    startTime,
		EndTime:      endTime,
		Description:  req.Description,
		OwnerID:      req.OwnerID,
		NotifyBefore: notifyBefore,
	}, nil
}

func (e *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	params, err := e.prepareForCreateOrUpdate(w, r)
	if err != nil {
		return
	}

	_, err = e.app.CreateEvent(r.Context(), *params)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExists) {
			RespondWithJSON(w, http.StatusBadRequest, "Event already exists")
			return
		}

		RespondWithJSON(w, http.StatusInternalServerError, Error("Failed to create event"))
		return
	}

	RespondWithJSON(w, http.StatusCreated, OK())
}

func (e *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("id")
	if eventID == "" {
		RespondWithJSON(w, http.StatusBadRequest, Error("Event ID is required"))
		return
	}

	param, err := e.prepareForCreateOrUpdate(w, r)
	if err != nil {
		return
	}

	event := storage.Event{
		ID:           eventID,
		Title:        param.Title,
		StartTime:    param.StartTime,
		EndTime:      param.EndTime,
		Description:  param.Description,
		OwnerID:      param.OwnerID,
		NotifyBefore: param.NotifyBefore,
	}

	err = e.app.UpdateEvent(r.Context(), event)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExists) {
			RespondWithJSON(w, http.StatusBadRequest, "Event already exists")
			return
		}

		RespondWithJSON(w, http.StatusInternalServerError, Error("Failed to create event"))
		return
	}

	RespondWithJSON(w, http.StatusCreated, OK())
}

func (e *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("id")
	if eventID == "" {
		RespondWithJSON(w, http.StatusBadRequest, Error("Event ID is required"))
		return
	}

	err := e.app.DeleteEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			RespondWithJSON(w, http.StatusNotFound, Error("Event not found"))
			return
		}

		RespondWithJSON(w, http.StatusInternalServerError, Error("Failed to delete event"))
		return
	}

	RespondWithJSON(w, http.StatusOK, OK())
}

func (e *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("id")
	if eventID == "" {
		RespondWithJSON(w, http.StatusBadRequest, Error("Event ID is required"))
		return
	}

	event, err := e.app.GetEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			RespondWithJSON(w, http.StatusNotFound, Error("Event not found"))
			return
		}

		RespondWithJSON(w, http.StatusInternalServerError, Error("Failed to get event"))
		return
	}

	RespondWithJSON(w, http.StatusOK, event)
}

func (e *EventHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := e.app.GetAllEvents(r.Context())
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, Error("Failed to get events"))
		return
	}

	RespondWithJSON(w, http.StatusOK, events)
}
