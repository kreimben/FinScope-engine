package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL       string
	SupabaseAnonKey   string
	SupabaseJWTSecret string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		SupabaseURL:       os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:   os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseJWTSecret: os.Getenv("SUPABASE_JWT_SECRET"),
	}
}
