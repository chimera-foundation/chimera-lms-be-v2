package response

import (
	"encoding/json"
	"net/http"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/dto"
)

func base(w http.ResponseWriter, code int, status string, data any, errs any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := dto.WebResponse{
		Code:   code,
		Status: status,
		Data:   data,
		Errors: errs,
	}

	json.NewEncoder(w).Encode(res)
}

// --- Success Helpers ---

func OK(w http.ResponseWriter, data any) {
	base(w, http.StatusOK, "OK", data, nil)
}

func Created(w http.ResponseWriter, data any) {
	base(w, http.StatusCreated, "CREATED", data, nil)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// --- Client Error Helpers ---

func BadRequest(w http.ResponseWriter, message string) {
	base(w, http.StatusBadRequest, "BAD_REQUEST", nil, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	base(w, http.StatusUnauthorized, "UNAUTHORIZED", nil, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	base(w, http.StatusForbidden, "FORBIDDEN", nil, message)
}

func NotFound(w http.ResponseWriter, message string) {
	base(w, http.StatusNotFound, "NOT_FOUND", nil, message)
}

func UnprocessableEntity(w http.ResponseWriter, message string) {
	base(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", nil, message)
}

// --- Server Error Helpers ---

func InternalServerError(w http.ResponseWriter, message string) {
	base(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", nil, message)
}