package grpchandler

import (
	"context"
	"errors"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/helpers"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	pb "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/pb/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventHandler struct {
	pb.UnimplementedEventsServer
	app app.Application
}

func NewEventHandler(app app.Application) *EventHandler {
	return &EventHandler{app: app}
}

func (h *EventHandler) Create(ctx context.Context, req *pb.CreateOrUpdateEventRequest) (*pb.EmptyResponse, error) {
	param, err := createOrUpdateRequestToStorageParams(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = h.app.CreateEvent(ctx, *param)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (h *EventHandler) Get(ctx context.Context, req *pb.GetEventRequest) (*pb.Event, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Event ID is required")
	}

	event, err := h.app.GetEvent(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return eventToProto(*event), nil
}

func (h *EventHandler) Update(ctx context.Context, req *pb.CreateOrUpdateEventRequest) (*pb.EmptyResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Event ID is required")
	}

	param, err := createOrUpdateRequestToStorageParams(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event := storage.Event{
		ID:           req.GetId(),
		Title:        param.Title,
		StartTime:    param.StartTime,
		EndTime:      param.EndTime,
		Description:  param.Description,
		OwnerID:      param.OwnerID,
		NotifyBefore: param.NotifyBefore,
	}

	err = h.app.UpdateEvent(ctx, event)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (h *EventHandler) Delete(ctx context.Context, req *pb.DeleteEventRequest) (*pb.EmptyResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Event ID is required")
	}

	err := h.app.DeleteEvent(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, _ *pb.EmptyRequest) (*pb.EventListResponse, error) {
	events, err := h.app.GetAllEvents(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	eventList := make([]*pb.Event, len(events))
	for i, event := range events {
		eventList[i] = eventToProto(event)
	}

	return &pb.EventListResponse{Events: eventList}, nil
}

func (h *EventHandler) ListDayEvents(ctx context.Context, req *pb.DateRequest) (*pb.EventListResponse, error) {
	events, err := h.app.GetEventsForDay(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return eventsToResponse(events), nil
}

func (h *EventHandler) ListWeekEvents(ctx context.Context, req *pb.DateRequest) (*pb.EventListResponse, error) {
	events, err := h.app.GetEventsForWeek(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return eventsToResponse(events), nil
}

func (h *EventHandler) ListMonthEvents(ctx context.Context, req *pb.DateRequest) (*pb.EventListResponse, error) {
	events, err := h.app.GetEventsForMonth(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return eventsToResponse(events), nil
}

func eventsToResponse(events []storage.Event) *pb.EventListResponse {
	eventList := make([]*pb.Event, len(events))
	for i, event := range events {
		eventList[i] = eventToProto(event)
	}
	return &pb.EventListResponse{Events: eventList}
}

func createOrUpdateRequestToStorageParams(
	req *pb.CreateOrUpdateEventRequest,
) (*storage.CreateOrUpdateEventParams, error) {
	if req.GetStartTime() == nil || req.GetEndTime() == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end time must be provided")
	}

	if !helpers.IsValidUUID(req.GetOwnerId()) {
		return nil, status.Error(codes.InvalidArgument, "owner_id must be uuid")
	}

	var notifyBefore time.Duration
	if req.GetNotifyBefore() != nil {
		notifyBefore = req.GetNotifyBefore().AsDuration()
	}

	return &storage.CreateOrUpdateEventParams{
		Title:        req.GetTitle(),
		StartTime:    req.GetStartTime().AsTime(),
		EndTime:      req.GetEndTime().AsTime(),
		Description:  req.Description,
		OwnerID:      req.GetOwnerId(),
		NotifyBefore: &notifyBefore,
	}, nil
}

func eventToProto(e storage.Event) *pb.Event {
	eventProto := &pb.Event{
		Id:        e.ID,
		Title:     e.Title,
		StartTime: timestamppb.New(e.StartTime),
		EndTime:   timestamppb.New(e.EndTime),
		OwnerId:   e.OwnerID,
	}

	if e.NotifyBefore != nil {
		eventProto.NotifyBefore = durationpb.New(*e.NotifyBefore)
	}

	return eventProto
}
