package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	JWTSecret   string
	FrontendURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPass := os.Getenv("DB_PASSWORD")

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pkbm_digital_darel_iman"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dar-el-iman-secret-key-2027-2032"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return &Config{
		Port:        port,
		DBUser:      dbUser,
		DBPassword:  dbPass,
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBName:      dbName,
		JWTSecret:   jwtSecret,
		FrontendURL: frontendURL,
	}
}
