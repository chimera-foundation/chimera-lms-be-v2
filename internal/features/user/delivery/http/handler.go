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
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	authService service.Auth
	log         *logrus.Logger
}

func (h *UserHandler) PublicRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register/student", h.Register)
	r.Post("/login", h.Login)

	return r
}

func (h *UserHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/logout", h.Logout)
	r.Get("/me", h.Me)

	return r
}

func NewUserHandler(authService service.Auth, log *logrus.Logger) *UserHandler {
	return &UserHandler{authService: authService, log: log}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("invalid request payload for register")
		u.BadRequest(w, "Invalid request payload")
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		h.log.WithError(err).WithField("org_id", req.OrganizationID).Error("invalid organization ID")
		u.InternalServerError(w, err.Error())
		return
	}
	user, err := h.authService.RegisterStudent(
		r.Context(),
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		orgID,
	)

	if err != nil {
		h.log.WithError(err).WithField("email", req.Email).Error("failed to register student")
		u.InternalServerError(w, err.Error())
		return
	}

	response := dto.RegisterResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	u.Created(w, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("invalid request payload for login")
		u.BadRequest(w, "Invalid request payload")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.WithField("email", req.Email).Warn("login failed")
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
		h.log.WithError(err).Error("logout failed")
		u.InternalServerError(w, "Failed to logout")
		return
	}

	u.OK(w, map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.Me(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to fetch user information")
		u.InternalServerError(w, "Failed to fetch user information")
		return
	}

	response := dto.MeResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     user.RolesStr(),
	}
	u.OK(w, response)
}
