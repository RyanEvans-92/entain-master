package service

import (
	"golang.org/x/net/context"

	"sports/proto/sports"

	"sports/db"
)

type Events interface {
	// ListEvents will return a collection of events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// eventsService implements the Events interface.
type eventsService struct {
	eventsRepo db.EventsRepo
}

// NewEventsService instantiates and returns a new eventsService.
func NewEventsService(eventsRepo db.EventsRepo) Events {
	return &eventsService{eventsRepo}
}

func (s *eventsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	events, err := s.eventsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sports.ListEventsResponse{Events: events}, nil
}
