package config

import (
	"os"
)

type Config struct {
	SupabaseURL        string
	SupabaseKey        string
	SupabaseProjectRef string
	BaseURL            string
	QRSize             int
}

func Load() *Config {
	return &Config{
		SupabaseURL:        getEnv("SUPABASE_URL", "https://your-project.supabase.co"),
		SupabaseKey:        getEnv("SUPABASE_KEY", "your-anon-key"),
		SupabaseProjectRef: getEnv("SUPABASE_PROJECT_REF", "your-project-ref"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		QRSize:             getEnvAsInt("QR_SIZE", 256),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}
