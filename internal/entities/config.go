package entities

import "encoding/json"

type Config struct {
	Token string `json:"token"`
}

func ParseConfig(rawConfig []byte) (Config, error) {
	config := Config{}
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
