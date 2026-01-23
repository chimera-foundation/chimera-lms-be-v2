package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/service"
	sdto "github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	authService service.Auth
}

func NewUserHandler(authService service.Auth) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) respondWithError(w http.ResponseWriter, code int, status string, message string) {
    response := sdto.WebResponse{
        Code:   code,
        Status: status,
        Errors: message,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) respondWithJSON(w http.ResponseWriter, code int, status string, payload any) {
    response := sdto.WebResponse{
        Code:   code,
        Status: status,
        Data:   payload,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	
	return r
}

func (h *UserHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()
	
    r.Post("/logout", h.Logout)
	
	return r
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request payload")
        return
    }

    orgID, err := uuid.Parse(req.OrganizationID)
    if err != nil {
        h.respondWithError(w, http.StatusInternalServerError,"INTERNAL_SERVER_ERROR", err.Error())
        return
    }
    user, err := h.authService.Register(
        r.Context(), 
        req.Email, 
        req.Password, 
        req.FirstName, 
        req.LastName, 
        orgID,
    )

    if err != nil {
        h.respondWithError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
        return
    }

    response := dto.RegisterResponse{
        Email: user.Email,
        FirstName: user.FirstName,
        LastName: user.LastName,
    }
    h.respondWithJSON(w, http.StatusCreated, "CREATED", response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    req := dto.LoginRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request payload")
        return
    }

    token, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        h.respondWithError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
        return
    }

    h.respondWithJSON(w, http.StatusOK, "OK",map[string]string{
        "access_token": token,
        "token_type":   "Bearer",
    })
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    err := h.authService.Logout(r.Context(), tokenString)
    if err != nil {
        h.respondWithError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to logout")
        return
    }

    h.respondWithJSON(w, http.StatusOK, "OK", map[string]string{
        "message": "Successfully logged out",
    })
}