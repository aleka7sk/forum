package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	Port        string
	Hash_salt   string
	Signing_key string
	Token_ttl   time.Duration
}

func InitConfig() (Config, error) {
	bytes, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
