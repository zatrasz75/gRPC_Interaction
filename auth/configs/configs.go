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
		ConnStr string `yaml:"conn-str" env:"DB_CONNECTION_STRING"`

		Host     string `yaml:"host" env:"HOST_DB"`
		User     string `yaml:"username" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		Url      string `yaml:"db-url" env:"URL_DB"`
		Name     string `yaml:"db-name" env:"POSTGRES_DB"`
		Port     string `yaml:"port" env:"PORT_DB"`

		PoolMax      int           `yaml:"pool-max" env:"PG_POOL_MAX" env-default:"2"`
		ConnAttempts int           `yaml:"conn-attempts" env:"PG_CONN_ATTEMPTS" env-default:"7"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"PG_TIMEOUT" env-default:"3s"`
	} `yaml:"database"`
	Token struct {
		Salt          string        `yaml:"secret-password-salt" env:"SECRET_PASSWORD_SALT"`
		SecretKeyHere string        `yaml:"secret-key-here" env:"SECRET_KEY_TOKEN"`
		Expiration    time.Duration `yaml:"expiration" env:"EXPIRATION_TIME"`
	} `yaml:"token"`
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
