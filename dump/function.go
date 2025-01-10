package dump

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arfan21/go-pgdump/config"
	"github.com/arfan21/go-pgdump/util"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type FunctionData struct {
	FunctionName string
	FunctionDef  string
}

func DumpFunctions(db *sql.DB, cfg config.Config) error {
	functions, err := getFunctionDefinitions(db)
	if err != nil {
		return fmt.Errorf("failed to get functions : %w", err)
	}

	if len(functions) == 0 {
		slog.Info("No functions found")
		return nil
	}

	slog.Info("Found functions for dumping", "count", len(functions))

	if err := dumpFunctionsToFile(cfg, functions); err != nil {
		return fmt.Errorf("failed to write function dump : %w", err)
	}

	return nil
}

func getFunctionDefinitions(db *sql.DB) ([]FunctionData, error) {
	query := `
		SELECT f.proname, pg_get_functiondef(f.oid)
		FROM pg_catalog.pg_proc f
		INNER JOIN pg_catalog.pg_namespace n ON (f.pronamespace = n.oid)
		WHERE n.nspname = 'public';
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error when getting function: %w", err)
	}
	defer rows.Close()
	var functions []FunctionData
	for rows.Next() {
		var functionData FunctionData
		if err := rows.Scan(&functionData.FunctionName, &functionData.FunctionDef); err != nil {
			return nil, fmt.Errorf("failed to scan function: %w", err)
		}
		functions = append(functions, functionData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during function scan: %w", err)
	}

	return functions, nil
}

func dumpFunctionsToFile(cfg config.Config, functions []FunctionData) error {
	errg, _ := errgroup.WithContext(context.Background())
	errg.SetLimit(5)
	for _, function := range functions {
		errg.Go(func() error {

			dir := "functions"
			if strings.HasPrefix((function.FunctionName), "sp") {
				dir = "procedures"
			}
			dir = filepath.Join(cfg.DumpDir, dir)
			err := util.CreateDirIfNotExist(dir)
			if err != nil {
				return fmt.Errorf("failed to create dir: %w", err)
			}

			filename := filepath.Join(dir, fmt.Sprintf("%s.sql", function.FunctionName))
			f, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer f.Close()

			_, err = f.WriteString(function.FunctionDef)
			if err != nil {
				return fmt.Errorf("failed to write function dump: %w", err)
			}

			slog.Info("Dumped function", dir, function.FunctionName)

			return nil
		})
	}

	err := errg.Wait()
	if err != nil {
		return fmt.Errorf("failed to dump functions: %w", err)
	}

	return nil
}
