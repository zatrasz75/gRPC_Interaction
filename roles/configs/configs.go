package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	GRPC struct {
		AddrHost string `yaml:"host" env:"APP_IP" env-default:"localhost"`
		AddrPort string `yaml:"port" env:"APP_PORT" env-default:"50051"`
	} `yaml:"GRPC"`
	PostgreSQL struct {
		ConnStr string `env:"DB_CONNECTION_STRING" env-description:"db string"`

		Host     string `yaml:"host" env:"HOST_DB" env-description:"db host"`
		User     string `yaml:"username" env:"POSTGRES_USER" env-description:"db username"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-description:"db password"`
		Url      string `yaml:"db-url" env:"URL_DB" env-description:"db url"`
		Name     string `yaml:"db-name" env:"POSTGRES_DB" env-description:"db name"`
		Port     string `yaml:"port" env:"PORT_DB" env-description:"db port"`

		PoolMax      int           `yaml:"pool-max" env:"PG_POOL_MAX" env-description:"db PoolMax" env-default:"2"`
		ConnAttempts int           `yaml:"conn-attempts" env:"PG_CONN_ATTEMPTS" env-description:"db ConnAttempts" env-default:"5"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"PG_TIMEOUT" env-description:"db ConnTimeout" env-default:"2s"`
	} `yaml:"database"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	cfg.PostgreSQL.ConnStr = initDB(cfg)

	return &cfg, nil
}

func initDB(cfg Config) string {
	if cfg.PostgreSQL.ConnStr != "" {
		return cfg.PostgreSQL.ConnStr
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Url,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Name,
	)
}