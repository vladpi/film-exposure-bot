package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBDSN    string `envconfig:"DB_DSN"`
	BotToken string `envconfig:"BOT_TOKEN"`
}

func LoadFromDotenv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
		// log.Fatal("Error loading .env file")
	}

	var c Config
	err = envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
