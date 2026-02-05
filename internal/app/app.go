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

	assessmentHttp "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/delivery/http"
	assessmentPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/repository/postgres"
	assessmentService "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/service"
	cohortPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/repository/postgres"
	enrollmentPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/repository/postgres"
	eventHttp "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/delivery/http"
	eventPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/repository/postgres"
	eventService "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/service"
	orgPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	sectionPostgres "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/repository/postgres"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/middleware"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
)

type BootstrapConfig struct {
	DB     *sql.DB
	Router *chi.Mux
	Log    *logrus.Logger
	Config *viper.Viper
	Redis  *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Setup Repositories
	userRepo := postgres.NewUserRepo(config.DB)
	roleRepo := postgres.NewRoleRepository(config.DB)

	// Event Dependencies
	eventRepo := eventPostgres.NewEventRepository(config.DB)
	orgRepo := orgPostgres.NewOrganizationRepo(config.DB)
	enrollmentRepo := enrollmentPostgres.NewEnrollmentRepository(config.DB)
	cohortRepo := cohortPostgres.NewCohortRepository(config.DB)
	sectionRepo := sectionPostgres.NewSectionRepository(config.DB)

	// Assessment Dependencies
	assessmentRepo := assessmentPostgres.NewAssessmentRepoPostgres(config.DB)

	secret := config.Config.GetString("JWT_SECRET_KEY")
	expiryMinutes := config.Config.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")
	if expiryMinutes == 0 {
		expiryMinutes = 60
	}
	expiryDuration := time.Duration(expiryMinutes) * time.Minute

	tokenProvider := auth.NewJWTProvider(secret, expiryDuration, config.Redis)

	// 2. Setup Services/UseCases
	authService := service.NewAuthService(userRepo, roleRepo, tokenProvider)

	eventService := eventService.NewEventService(
		eventRepo,
		orgRepo,
		enrollmentRepo,
		cohortRepo,
		sectionRepo,
		config.Redis,
	)

	assessmentSvc := assessmentService.NewAssessmentService(assessmentRepo)

	// 3. Setup Controllers/Handlers
	userHandler := http.NewUserHandler(authService)
	eventHandler := eventHttp.NewEventHandler(eventService)
	assessmentHandler := assessmentHttp.NewAssessmentHandler(assessmentSvc)

	// 4. Setup Routes
	config.Router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Mount("/auth", userHandler.PublicRoutes())
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(tokenProvider))

			r.Mount("/users", userHandler.ProtectedRoutes())
			r.Mount("/events", eventHandler.ProtectedRoutes())
			r.Mount("/assessments", assessmentHandler.ProtectedRoutes())
		})
	})
}
