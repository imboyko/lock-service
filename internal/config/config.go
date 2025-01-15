package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load() (Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}

type Config struct {
	JwtSecret string `env:"JWT_SECRET" env-required:"true"`
	Redis
}

type Redis struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Db       int    `env:"REDIS_DB"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`
}

func (r *Redis) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
