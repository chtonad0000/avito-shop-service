package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig(flag bool) (*Config, error) {

	var dbURL string
	if flag {
		dbURL = os.Getenv("TEST_DATABASE_URL")
	} else {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	return &Config{
		DatabaseURL: dbURL,
	}, nil
}
