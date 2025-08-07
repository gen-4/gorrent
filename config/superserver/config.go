package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	DEV  = "dev"
	TEST = "test"
	PRO  = "pro"
)

type SuperserverConfig struct {
	TorrentsFolder  string `json:"torrents_folder"`
	LogFile         string `json:"log_file"`
	Peers           []string
	Environment     string
	PeerUrlTemplate string
}

var Configuration SuperserverConfig = SuperserverConfig{
	TorrentsFolder:  "",
	LogFile:         "server.log",
	Peers:           []string{},
	Environment:     PRO,
	PeerUrlTemplate: "%s://%s:%s/gorrent/%s",
}

var fileDescriptor *os.File

func getEnv() string {
	var err error
	environment := DEV

	if flag.Lookup("test.v") == nil {
		err = godotenv.Load()
	} else {
		envFileError := godotenv.Load(".test.env")
		if envFileError != nil {
			err = godotenv.Load("../.test.env")
		}
	}

	if err != nil {
		slog.Warn("Unable to read .env file")
	} else {
		environment = os.Getenv("ENVIRONMENT")
	}

	return environment
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
	var environment string = getEnv()
	Configuration.Environment = environment
	f, err := os.OpenFile(Configuration.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error opening log file", "error", err.Error())
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
	slog.SetDefault(logger)

	switch environment {
	case DEV:
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		f.Close()
		Configuration.PeerUrlTemplate = fmt.Sprintf(Configuration.PeerUrlTemplate, "http", "%s", "5050", "%s")

	case TEST:
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		fileDescriptor = f
		Configuration.PeerUrlTemplate = fmt.Sprintf(Configuration.PeerUrlTemplate, "http", "%s", "5050", "%s")

	case PRO:
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
		slog.SetDefault(logger)
		fileDescriptor = f
		Configuration.PeerUrlTemplate = fmt.Sprintf(Configuration.PeerUrlTemplate, "http", "%s", "5050", "%s")
	}

	slog.Info(fmt.Sprintf("Running in %s environment", environment))
}

func CloseConfig() {
	if fileDescriptor != nil {
		fileDescriptor.Close()
	}

}
