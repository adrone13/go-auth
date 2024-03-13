package config

import (
	"fmt"
	"github.com/adrone13/goenvconfig"
	"github.com/joho/godotenv"
	"log"
)

var Values *Config

type Config struct {
	Port      int    `env:"PORT"`
	JwtSecret string `env:"JWT_SECRET"`
	JwtTtl    int    `env:"JWT_TTL"`

	DbHost     string `env:"DB_HOST"`
	DbName     string `env:"DB_NAME"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbPort     int    `env:"DB_PORT"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	Values = new(Config)

	err = goenvconfig.Load(Values)
	if err != nil {
		panic(fmt.Sprintf("Failed to load env config. Error: %v", err))
	}
}
