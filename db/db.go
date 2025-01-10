package db

import (
	"database/sql"
	"fmt"
	"github.com/arfan21/go-pgdump/config"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func NewDB(cfg config.Config) *sql.DB {
	host, port, user, password, dbname, err := cfg.ParseDatabaseURL()
	if err != nil {
		slog.Error("Failed to parse database URL", "error", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to database", "host", host, "port", port, "dbname", dbname)
	return db
}
