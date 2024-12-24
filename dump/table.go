package dump

import (
	"context"
	"database/sql"
	"fmt"
	"go-pgdump/config"
	"go-pgdump/util"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func DumpTables(db *sql.DB, cfg config.Config) error {
	tables, err := getTables(db)
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		slog.Info("No tables found")
		return nil
	}

	slog.Info("Found tables for dumping", "count", len(tables))
	errg, _ := errgroup.WithContext(context.Background())
	errg.SetLimit(5) // set max concurrent dumps
	for _, tableName := range tables {
		errg.Go(func() error {
			err := dumpTableToFile(cfg, tableName)
			if err != nil {
				slog.Error("Failed to dump table", "table", tableName, "error", err)
			}

			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return fmt.Errorf("failed to dump tables: %w", err)
	}

	return nil
}

func getTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT 
			ist.table_name
		FROM
			information_schema.tables ist
		WHERE
			ist.table_type = 'BASE TABLE' 
			AND ist.table_schema = 'public'
		ORDER BY ist.table_name;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableData string
		if err := rows.Scan(&tableData); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, tableData)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during table scan: %w", err)
	}

	return tables, nil
}

func dumpTableToFile(cfg config.Config, tableName string) error {
	dir := cfg.DumpDir + "/tables"
	err := util.CreateDirIfNotExist(dir)
	if err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}
	filename := filepath.Join(dir, fmt.Sprintf("%s.sql", tableName))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	host, port, user, password, dbname, err := cfg.ParseDatabaseURL()
	if err != nil {
		return fmt.Errorf("failed to parse db url: %w", err)
	}

	os.Setenv("PGPASSWORD", password)

	// table dump
	cmd := exec.Command("pg_dump",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbname,
		"-t", tableName,
		"--schema-only",
		"--no-owner",
		"--no-privileges",
	)
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to dump table %s: %w", tableName, err)
	}
	slog.Info("Dumped table to", "filename", filename)
	return nil
}
