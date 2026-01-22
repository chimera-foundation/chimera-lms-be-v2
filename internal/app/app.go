package app

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/delivery/http"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/service"
)

type BootstrapConfig struct {
	DB     *sql.DB
	Router *chi.Mux     
	Log    *logrus.Logger
	Config *viper.Viper
	TokenProvider  service.TokenProvider
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Setup Repositories
	userRepo := postgres.NewUserRepo(config.DB)
	// tokenService := auth.NewJWTProvider()

	// 2. Setup Services/UseCases
	authService := service.NewAuthService(userRepo, config.TokenProvider) // nil is placeholder for TokenProvider

	// 3. Setup Controllers/Handlers
	userHandler := http.NewUserHandler(authService)

	// 4. Setup Routes (Standard library way)
	config.Router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userHandler.Routes())
	})
}