package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBSource             string
	JWTSecret            string
	DefaultAdminUser     string
	DefaultAdminPassword string
	AppPort              string
	RecaptchaSecretKey   string
}

// LoadConfig loads configuration from .env file
func LoadConfig(path string) (*Config, error) {
	// if err := godotenv.Load(path); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	_ = godotenv.Load(path)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"))

	return &Config{
		DBSource:             dsn,
		JWTSecret:            os.Getenv("JWT_SECRET"),
		DefaultAdminUser:     os.Getenv("DEFAULT_ADMIN_USER"),
		DefaultAdminPassword: os.Getenv("DEFAULT_ADMIN_PASSWORD"),
		AppPort:              os.Getenv("APP_PORT"),
		RecaptchaSecretKey:   os.Getenv("RECAPTCHA_SECRET_KEY"),
	}, nil
}
