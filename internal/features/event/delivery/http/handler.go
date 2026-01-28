package http

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/service"
	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	eventService service.EventService
}

func (h *EventHandler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	
	return r
}

func (h *EventHandler) PrivateRoutes() chi.Router {
	r := chi.NewRouter()
	
	return r
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}