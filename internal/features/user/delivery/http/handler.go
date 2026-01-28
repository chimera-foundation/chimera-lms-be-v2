package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/service"
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
    authService service.Auth
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
    r.Get("/me", h.Me)
    
    return r
}

func NewUserHandler(authService service.Auth) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        u.BadRequest(w, "Invalid request payload")
        return
    }

    orgID, err := uuid.Parse(req.OrganizationID)
    if err != nil {
        u.InternalServerError(w, err.Error())
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
        u.InternalServerError(w, err.Error())
        return
    }

    response := dto.RegisterResponse{
        Email: user.Email,
        FirstName: user.FirstName,
        LastName: user.LastName,
    }

    u.Created(w, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    req := dto.LoginRequest{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        u.BadRequest(w, "Invalid request payload")
        return
    }

    token, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        u.InternalServerError(w, err.Error())
        return
    }

    u.OK(w, map[string]string{
        "access_token": token,
        "token_type":   "Bearer",
    })
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    err := h.authService.Logout(r.Context(), tokenString)
    if err != nil {
        u.InternalServerError(w, "Failed to logout")
        return
    }

    u.OK(w, map[string]string{
        "message": "Successfully logged out",
    })
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    user, err := h.authService.Me(r.Context(), tokenString)
    if err != nil {
        u.InternalServerError(w, "Failed to fetch user information")
        return
    }

    response := dto.MeResponse{
        Email: user.Email,
        FirstName: user.FirstName,
        LastName: user.LastName,
    }
    u.OK(w, response)
}