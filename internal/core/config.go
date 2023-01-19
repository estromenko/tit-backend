package core

import "github.com/tutorin-tech/tit-backend/internal/utils"

const (
	defaultPort           = 3000
	defaultPgPort         = 5432
	defaultJWTExpireHours = 24 * 3
)

type Config struct {
	Debug          bool
	Port           int
	PgHost         string
	PgPort         int
	PgName         string
	PgUser         string
	PgPassword     string
	SecretKey      string
	JWTExpireHours int
}

func NewConfig() *Config {
	return &Config{
		Debug:          utils.GetEnvOrDefault("DEBUG", "false") == "true",
		Port:           utils.GetEnvIntOrDefault("PORT", defaultPort),
		PgHost:         utils.GetEnvOrDefault("PG_HOST", "localhost"),
		PgPort:         utils.GetEnvIntOrDefault("PG_PORT", defaultPgPort),
		PgName:         utils.GetEnvOrDefault("PG_NAME", "tutorintech"),
		PgUser:         utils.GetEnvOrDefault("PG_USER", "postgres"),
		PgPassword:     utils.GetEnvOrDefault("PG_PASSWORD", "secret"),
		SecretKey:      utils.GetEnv("SECRET_KEY"),
		JWTExpireHours: utils.GetEnvIntOrDefault("JWT_EXPIRE_HOURS", defaultJWTExpireHours),
	}
}
