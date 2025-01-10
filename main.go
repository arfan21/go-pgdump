package main

import (
	"github.com/arfan21/go-pgdump/config"
	"github.com/arfan21/go-pgdump/db"
	"github.com/arfan21/go-pgdump/dump"
	"github.com/arfan21/go-pgdump/util"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
	cfg := config.LoadConfig()

	err := util.CreateDirIfNotExist(cfg.DumpDir)
	if err != nil {
		slog.Error("Failed to create dump directory", "error", err)
		os.Exit(1)
	}

	dbConn := db.NewDB(cfg)
	defer dbConn.Close()

	slog.Info("Start dumping tables")
	if err := dump.DumpTables(dbConn, cfg); err != nil {
		slog.Error("Failed to dump tables", "error", err)
	}
	slog.Info("Finish dumping tables")

	slog.Info("Start dumping functions")
	if err := dump.DumpFunctions(dbConn, cfg); err != nil {
		slog.Error("Failed to dump functions", "error", err)
	}
	slog.Info("Finish dumping functions")
}
