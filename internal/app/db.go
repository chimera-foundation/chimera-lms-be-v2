package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// Import the pgx driver for database/sql compatibility
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(v *viper.Viper, log *logrus.Logger) *sql.DB {
	dsn := v.GetString("DB_URI")
	if dsn == "" {
		log.Fatal("DB_URI is not set in environment variables")
	}

	// 2. Open the connection using the 'pgx' driver name
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	idleConn := v.GetInt("DB_POOL_IDLE")
	if idleConn == 0 { idleConn = 10 } 
	
	maxConn := v.GetInt("DB_POOL_MAX")
	if maxConn == 0 { maxConn = 100 } 
	
	maxLifetime := v.GetInt("DB_POOL_LIFETIME")
	if maxLifetime == 0 { maxLifetime = 300 } 

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Info("Database connection established successfully")
	return db
}