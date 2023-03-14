package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Config struct {
	Debug            bool   `envconfig:"debug"`
	Port             int    `envconfig:"port"`
	PostgresHost     string `envconfig:"postgres_host"`
	PostgresUser     string `envconfig:"postgres_user"`
	PostgresDB       string `envconfig:"postgres_db"`
	PostgresPort     int    `envconfig:"postgres_port"`
	PostgresPassword string `envconfig:"postgres_password"`
	Env              string `envconfig:"env"`
	RedisHost        string `envconfig:"redis_host"`
	RedisPort        string `envconfig:"redis_port"`
	RedisPassword    string `envconfig:"redis_password"`
}

func Load(filename string) (*Config, error) {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load(filename); err != nil {
			log.Printf("couldn't load env vars: %v", err)
		}
	}

	c := &Config{}
	err := envconfig.Process("meddle", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
