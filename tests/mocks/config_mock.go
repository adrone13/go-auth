package mocks

import "auth/internal/config"

var ConfigMock = &config.Config{
	Port:                    8080,
	JwtSecret:               "jwt_secret",
	JwtTtl:                  120,
	RefreshTokenAbsoluteTtl: 1200,
	DbHost:                  "127.0.0.1",
	DbName:                  "auth",
	DbUser:                  "db_user",
	DbPassword:              "db_pass",
	DbPort:                  5430,
}
