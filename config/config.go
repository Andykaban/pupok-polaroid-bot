package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type BotConfig struct {
	BotToken string `json:"bot_token"`
	BotProxyUrl string `json:"bot_proxy_url"`
	BotProxyLogin string `json:"bot_proxy_login"`
	BotProxyPassword string `json:"bot_proxy_password"`
	BotTempDir string `json:"bot_temp_dir"`
	FontPath string `json:"font_path"`
	BackgroundPath string `json:"background_path"`
}

func ParseBotConfig(configPath string) *BotConfig {
	var botConfig BotConfig
	raw, err := os.Open(configPath)
	defer raw.Close()
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(raw)
	err = json.Unmarshal(byteValue, &botConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &botConfig
}