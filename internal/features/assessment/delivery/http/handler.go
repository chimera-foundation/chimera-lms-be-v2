package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/service"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	response "github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

type AssessmentHandler struct {
	assessmentService service.AssessmentService
}

func NewAssessmentHandler(assessmentService service.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
	}
}

func (h *AssessmentHandler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	return r
}

func (h *AssessmentHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/student", h.GetStudentAssessments)

	return r
}

func (h *AssessmentHandler) GetStudentAssessments(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	// Parse query parameters
	filter := domain.StudentAssessmentFilter{
		Limit:  20, // default
		Offset: 0,
	}

	// Parse type filter
	if typeStr := r.URL.Query().Get("type"); typeStr != "" {
		assessmentType := domain.AssessmentType(typeStr)
		filter.Type = &assessmentType
	}

	// Parse date filters
	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		startTime, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			filter.StartDate = &startTime
		}
	}

	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		endTime, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			filter.EndDate = &endTime
		}
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	result, err := h.assessmentService.GetStudentAssessments(r.Context(), userID, filter)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, result)
}
