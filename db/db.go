package db

import (
	"chatapp/config"
	"database/sql"
	"fmt"

	"context"
	"time"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("postgres", config.PGUrl)
	if err != nil {
		panic(fmt.Errorf("error connecting to postgres %w", err))
	}

	pingCtx, workersCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer workersCancel()
	if err := db.PingContext(pingCtx); err != nil {
		panic(fmt.Errorf("error pinging postgres db: %w", err))
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetConnMaxIdleTime(15 * time.Minute)

	DB = db
}
