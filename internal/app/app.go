package app

import (
	"database/sql"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/delivery/http"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/service"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/middleware"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
)

type BootstrapConfig struct {
	DB     *sql.DB
	Router *chi.Mux     
	Log    *logrus.Logger
	Config *viper.Viper
	Redis *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Setup Repositories
	userRepo := postgres.NewUserRepo(config.DB)

	secret := config.Config.GetString("JWT_SECRET_KEY")
	expiryMinutes := config.Config.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")
	if expiryMinutes == 0 {
		expiryMinutes = 60 
	}
	expiryDuration := time.Duration(expiryMinutes) * time.Minute
	
	tokenProvider := auth.NewJWTProvider(secret, expiryDuration, config.Redis)

	// 2. Setup Services/UseCases
	authService := service.NewAuthService(userRepo, tokenProvider) 

	// 3. Setup Controllers/Handlers
	userHandler := http.NewUserHandler(authService)

	// 4. Setup Routes (Standard library way)
	config.Router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
            r.Mount("/auth", userHandler.PublicRoutes()) 
        })

		r.Group(func(r chi.Router) {
            r.Use(middleware.AuthMiddleware(tokenProvider))
            
            r.Mount("/users", userHandler.ProtectedRoutes())
        })
	})
}