package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port     string
	DBDriver string
	DBSource string
}

func Load() *Config {
	loadDotEnv(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}

	source := os.Getenv("DB_SOURCE")
	if source == "" {
		source = "clientstore"
	}

	return &Config{
		Port:     port,
		DBDriver: driver,
		DBSource: source,
	}
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}
