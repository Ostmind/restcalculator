package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Port    string        `json:"port"`
	EnvType string        `json:"envType"`
	Timeout time.Duration `json:"timeout"`
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
