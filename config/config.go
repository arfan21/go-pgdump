package config

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	DumpDir     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dumpDir := os.Getenv("DUMP_DIR")
	if dumpDir == "" {
		dumpDir = "db_dump"
	}

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DumpDir:     os.Getenv("DUMP_DIR"),
	}
}

func (c Config) ParseDatabaseURL() (string, string, string, string, string, error) {
	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		return "", "", "", "", "", err
	}

	password, _ := u.User.Password()
	queryParams := u.Query()

	sslMode := queryParams.Get("sslmode")
	if sslMode == "" {
		sslMode = "disable" // set default ssl mode
	}

	return u.Hostname(), u.Port(), u.User.Username(), password, u.Path[1:], nil
}
