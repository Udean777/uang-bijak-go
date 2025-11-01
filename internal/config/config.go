package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	AppPort         string
	JwtSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable tidak di-set!")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Nilai default jika APP_PORT tidak disetel
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable tidak di-set!")
	}

	accessTTL, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TTL_MINUTES"))
	if accessTTL == 0 {
		accessTTL = 15 // Default 15 menit
	}

	refreshTTL, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TTL_DAYS"))
	if refreshTTL == 0 {
		refreshTTL = 7 // Default 7 hari
	}

	return &Config{
		DatabaseURL:     dbURL,
		AppPort:         appPort,
		JwtSecret:       jwtSecret,
		AccessTokenTTL:  time.Minute * time.Duration(accessTTL),
		RefreshTokenTTL: time.Hour * 24 * time.Duration(refreshTTL),
	}
}
