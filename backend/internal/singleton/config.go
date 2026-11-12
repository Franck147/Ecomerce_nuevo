package singleton

import (
	"os"
	"sync"
)

// Config representa las variables de configuracion del sistema.
type Config struct {
	HTTPPort    string
	DatabaseURL string
	Currency    string
}

var (
	configInstance *Config
	configOnce     sync.Once
)

// GetConfig implementa el PATRON SINGLETON para la configuracion.
// Garantiza que la configuracion solo se cargue y valide una sola vez en memoria.
func GetConfig() *Config {
	configOnce.Do(func() {
		port := os.Getenv("HTTP_PORT")
		if port == "" {
			port = "8081" // Puerto 8081 para no colisionar con el proyecto padre
		}

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://ecommerce:ecommerce@localhost:5432/ecommerce_patrones?sslmode=disable"
		}

		currency := os.Getenv("CURRENCY")
		if currency == "" {
			currency = "PEN"
		}

		configInstance = &Config{
			HTTPPort:    port,
			DatabaseURL: dbURL,
			Currency:    currency,
		}
	})

	return configInstance
}
