package config

import (
	"github.com/adrone13/goenvconfig"
	"github.com/joho/godotenv"
	"log"
)

var Values *Config

type Config struct {
	Port      int    `env:"PORT"`
	JwtSecret string `env:"JWT_SECRET"`

	// Might be better to move these to DB settings
	JwtTtl                  int `env:"JWT_TTL"`                    // seconds
	RefreshTokenAbsoluteTtl int `env:"REFRESH_ABSOLUTE_TOKEN_TTL"` // seconds
	RefreshTokenIdleTtl     int `env:"REFRESH_IDLE_TOKEN_TTL"`     // seconds

	DbHost     string `env:"DB_HOST"`
	DbName     string `env:"DB_NAME"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbPort     int    `env:"DB_PORT"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env. Error:", err)
	}

	Values = new(Config)

	err = goenvconfig.Load(Values)
	if err != nil {
		log.Fatalf("Failed to load env config. Error: %v", err)
	}
}
