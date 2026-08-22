package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	Env                    string
	DatabaseURL            string
	SupabaseURL            string
	SupabaseJWTSecret      string
	InternalJobSecret      string
	H3CaptureResolution    int
	H3TileResolution       int
	MaxSpeedMPS            float64
	MaxAccuracyMeters      float64
	MinActivityDurationSec int
	AllowGuestAuth         bool
	TerritoryDecayDays     int
}

func Load() Config {
	return Config{
		Port:                   getEnv("PORT", "8080"),
		Env:                    getEnv("ENV", "development"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseJWTSecret:      os.Getenv("SUPABASE_JWT_SECRET"),
		InternalJobSecret:      os.Getenv("INTERNAL_JOB_SECRET"),
		H3CaptureResolution:    getEnvInt("H3_CAPTURE_RESOLUTION", 10),
		H3TileResolution:       getEnvInt("H3_TILE_RESOLUTION", 7),
		MaxSpeedMPS:            getEnvFloat("MAX_SPEED_MPS", 7.0),
		MaxAccuracyMeters:      getEnvFloat("MAX_ACCURACY_METERS", 40.0),
		MinActivityDurationSec: getEnvInt("MIN_ACTIVITY_DURATION_SEC", 60),
		AllowGuestAuth:         getEnvBool("ALLOW_GUEST_AUTH", true),
		TerritoryDecayDays:     getEnvInt("TERRITORY_DECAY_DAYS", 14),
	}
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		return v == "1" || strings.EqualFold(v, "true") || v == "yes"
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
