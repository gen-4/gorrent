package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type SuperserverConfig struct {
	TorrentsFolder string `json:"torrents_folder"`
}

var Configuration SuperserverConfig = SuperserverConfig{
	TorrentsFolder: "",
}

func loadFromConfigFile() {
	content, err := os.ReadFile("superserver_conf.json")
	if err != nil {
		slog.Error("Unable to read superserver config file", "error", err.Error())
		return
	}

	if err = json.Unmarshal(content, &Configuration); err != nil {
		slog.Error("Unable to unmarshal config", "error", err.Error())
		return
	}
}

func Config() {
	loadFromConfigFile()
}
