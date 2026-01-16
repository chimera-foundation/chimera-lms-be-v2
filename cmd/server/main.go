package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/platform"
	userAdapters "github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/adapters"
	userService "github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/service"
)

func main() {
	// 1. Database
	dsn := "postgres://chimera_admin:chimera@localhost:5432/lms?sslmode=disable"
	db := platform.NewPostgresDB(dsn)

	// 2. Setup Gin
	r := gin.Default()

	// 3. CORS Middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// 4. Dependency Injection
	userRepo := userAdapters.NewBunUserRepo(db)
	userSvc  := userService.NewUserService(userRepo, "your-super-secret-key")
	userHdl  := &userAdapters.UserHandler{Svc: userSvc}

	// 5. Routes
	api := r.Group("/api/v1")
	{
		userHdl.RegisterRoutes(api.Group("/users"))
	}

	r.Run(":8080")
}