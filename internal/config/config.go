package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClientID   string `json:"clientId"`
	Endpoint   string `json:"endpoint"`
	RootCAPath string `json:"rootCAPath"`
	CertPath   string `json:"certPath"`
	KeyPath    string `json:"keyPath"`
	PubTopic   string `json:"pubTopic"`
	SubTopic   string `json:"subTopic"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
