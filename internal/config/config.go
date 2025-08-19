package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Database DatabaseConfig
	Api      ApiConfig
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USERNAME" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `env:"POSTGRES_DATABASE" env-required:"true"`
}

type ApiConfig struct {
	Host string `env:"API_HOST" env-default:"localhost"`
	Port int    `env:"API_PORT" env-default:"8080"`
}

func MustInit() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			panic(err)
		}
	}
	return &cfg
}
