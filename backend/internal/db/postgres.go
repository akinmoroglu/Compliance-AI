package db

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
	migrate()
	log.Println("database connected")
}

func migrate() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS checks (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			status          TEXT NOT NULL DEFAULT 'complete',
			input_copy      TEXT NOT NULL,
			content_hash    TEXT,
			result          JSONB,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		ALTER TABLE checks ADD COLUMN IF NOT EXISTS content_hash TEXT;
	`)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
