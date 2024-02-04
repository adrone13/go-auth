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
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	fmt.Println("Initializing env config")

	Values = new(Config)

	err = goenvconfig.Load(Values)
	if err != nil {
		panic(fmt.Sprintf("Failed to load env config. Error: %v", err))
	}
}
