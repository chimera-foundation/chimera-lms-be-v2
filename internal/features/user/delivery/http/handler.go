package http

import (
	"encoding/json"
	"net/http"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	authService service.Auth
}

func NewUserHandler(authService service.Auth) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) respondWithError(w http.ResponseWriter, code int, message string) {
    h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *UserHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	
	r.Post("/register", h.Register)
	r.Get("/login", h.Login)
	
	return r
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    user, err := h.authService.Register(
        r.Context(), 
        req.Email, 
        req.Password, 
        req.FirstName, 
        req.LastName, 
        uuid.MustParse(req.OrganizationID),
    )

    if err != nil {
        h.respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    req := dto.LoginRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    token, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        h.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    h.respondWithJSON(w, http.StatusOK, map[string]string{
        "access_token": token,
        "token_type":   "Bearer",
    })
}