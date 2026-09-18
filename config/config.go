package config

import "os"

type Config struct {
	Port              string
	CatalogServiceURL string
	OrderServiceURL   string
	WeatherServiceURL string
	JWTSecret         string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	if catalogURL == "" {
		catalogURL = "http://localhost:8081"
	}

	orderURL := os.Getenv("ORDER_SERVICE_URL")
	if orderURL == "" {
		orderURL = "http://localhost:8082"
	}

	weatherURL := os.Getenv("WEATHER_SERVICE_URL")
	if weatherURL == "" {
		weatherURL = "http://localhost:8083"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "agroconnect_super_secure_jwt_secret_key_2026"
	}

	return &Config{
		Port:              port,
		CatalogServiceURL: catalogURL,
		OrderServiceURL:   orderURL,
		WeatherServiceURL: weatherURL,
		JWTSecret:         jwtSecret,
	}
}
