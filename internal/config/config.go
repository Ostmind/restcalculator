package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Port            int           `json:"port"`
	EnvType         string        `json:"env_type"`
	ShutdownTimeout time.Duration `json:"timeout_seconds"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
