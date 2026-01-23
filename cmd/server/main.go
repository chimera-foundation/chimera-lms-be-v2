package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
)

func main() {
	v := app.NewViper()

	log := app.NewLogger(v)

	redis := app.NewRedis(v, log)

	db := app.NewDatabase(v, log)
	defer db.Close() 

	r := app.NewRouter()

	app.Bootstrap(&app.BootstrapConfig{
		DB:            db,
		Router:        r,
		Log:           log,
		Config:        v,
		Redis: redis,
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