package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/delivery/dto"
	e "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/service"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	response "github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type EventHandler struct {
	eventService service.EventService
	log          *logrus.Logger
}

func NewEventHandler(eventService service.EventService, log *logrus.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		log:          log,
	}
}

func (h *EventHandler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	return r
}

func (h *EventHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateEvent)
	r.Get("/calendar", h.GetCalendar)
	r.Get("/sections/{sectionID}/schedule", h.GetSectionSchedule)
	r.Get("/announcements", h.GetAnnouncements)
	r.Get("/", h.GetEvents) // General list, maybe for admin or global view
	return r
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("invalid request payload for create event")
		response.BadRequest(w, "Invalid request payload")
		return
	}

	orgID, ok := auth.GetOrgID(r.Context())
	if !ok {
		h.log.Warn("organization not found in context for create event")
		response.Unauthorized(w, "Organization not found in context")
		return
	}

	// Map DTO to Domain
	event := &e.Event{
		OrganizationID: orgID,
		Title:          req.Title,
		Description:    req.Description,
		Location:       req.Location,
		EventType:      e.EventType(req.EventType),
		Color:          req.Color,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		IsAllDay:       req.IsAllDay,
		RecurrenceRule: req.RecurrenceRule,
		Scope:          e.EventScope(req.Scope),
		CohortID:       req.CohortID,
		SectionID:      req.SectionID,
		UserID:         req.UserID,
		ImageURL:       req.ImageURL,
	}

	createdEvent, err := h.eventService.CreateEvent(r.Context(), event)
	if err != nil {
		h.log.WithError(err).Error("failed to create event")
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, createdEvent)
}

func (h *EventHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		response.BadRequest(w, "start and end query parameters are required")
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		response.BadRequest(w, "Invalid start time format (RFC3339 required)")
		return
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		response.BadRequest(w, "Invalid end time format (RFC3339 required)")
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	events, err := h.eventService.GetCalendarForUser(r.Context(), userID, startTime, endTime)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("failed to get calendar for user")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, events)
}

func (h *EventHandler) GetSectionSchedule(w http.ResponseWriter, r *http.Request) {
	sectionIDStr := chi.URLParam(r, "sectionID")
	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid section ID")
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		response.BadRequest(w, "start and end query parameters are required")
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		response.BadRequest(w, "Invalid start time format")
		return
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		response.BadRequest(w, "Invalid end time format")
		return
	}

	events, err := h.eventService.GetSectionSchedule(r.Context(), sectionID, startTime, endTime)
	if err != nil {
		h.log.WithError(err).WithField("section_id", sectionID).Error("failed to get section schedule")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, events)
}

func (h *EventHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	orgID, ok := auth.GetOrgID(r.Context())
	if !ok {
		response.Unauthorized(w, "Organization not found")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var startTime, endTime time.Time
	if startStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endStr)
	}

	events, err := h.eventService.GetAnnouncements(r.Context(), orgID, startTime, endTime, limit, offset)
	if err != nil {
		h.log.WithError(err).Error("failed to get announcements")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, events)
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	orgID, ok := auth.GetOrgID(r.Context())
	if !ok {
		response.Unauthorized(w, "Organization not found")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var startTime, endTime time.Time
	if startStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endStr)
	}

	events, err := h.eventService.GetEvents(r.Context(), orgID, startTime, endTime, limit, offset)
	if err != nil {
		h.log.WithError(err).Error("failed to get events")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, events)
}
