package config

import "os"

type Config struct {
	Port     string
	DBPath   string
	LogLevel string
	GinMode  string
}

func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "todo.db"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}

	return Config{Port: port, DBPath: dbPath, LogLevel: logLevel, GinMode: ginMode}
}
