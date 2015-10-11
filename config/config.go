package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
    BotApi string `json:"bot-api"`
    Database string `json:"database"`
}

func ParseConfig() *Config {
    data, err := ioutil.ReadFile("config/bot_config.json")

    if err != nil {
        panic("Error open bot_config.json")
    }

    configStruct := &Config{}

    json.Unmarshal(data, configStruct)

    return configStruct
}