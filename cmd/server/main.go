package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	v := app.NewViper()

	log := app.NewLogger(v)

	redis := app.NewRedis(v, log)

	db := app.NewDatabase(v, log)
	defer db.Close() 

	secret := v.GetString("JWT_SECRET_KEY")
	expiryMinutes := v.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")
	if expiryMinutes == 0 {
		expiryMinutes = 60 
	}
	expiryDuration := time.Duration(expiryMinutes) * time.Minute
	
	tokenProvider := auth.NewJWTProvider(secret, expiryDuration, redis)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) 
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	app.Bootstrap(&app.BootstrapConfig{
		DB:            db,
		Router:        r,
		Log:           log,
		Config:        v,
		TokenProvider: tokenProvider,
	})

	port := v.GetString("PORT")
	if port == "" {
		port = "8000"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Infof("LMS Server starting on %s", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}