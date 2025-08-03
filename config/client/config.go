package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const (
	DEV  = "dev"
	TEST = "test"
	PRO  = "pro"
)

type SuperserverConfig struct {
	LogFile                string   `json:"log_file"`
	Superservers           []string `json:"superservers"`
	Env                    string
	SuperserverUrlTemplate string
}

var Configuration SuperserverConfig = SuperserverConfig{
	LogFile:                "gorrent.log",
	Superservers:           []string{},
	Env:                    "",
	SuperserverUrlTemplate: "%s://%s:%s/gorrent/%s",
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
	content, err := os.ReadFile("gorrent_conf.json")
	if err != nil {
		slog.Error("Unable to read config file", "error", err.Error())
		return
	}

	if err = json.Unmarshal(content, &Configuration); err != nil {
		slog.Error("Unable to unmarshal config", "error", err.Error())
		return
	}
}

func Config() {
	var environment string = getEnv()
	Configuration.Env = environment
	loadFromConfigFile()
	logFile := Configuration.LogFile
	f, err := tea.LogToFile(logFile, "debug")
	if err != nil {
		slog.Error("Error opening log file", "error", err.Error())
		os.Exit(1)
	}
	fileDescriptor = f

	logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
	slog.SetDefault(logger)

	switch environment {
	case DEV:
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		Configuration.SuperserverUrlTemplate = fmt.Sprintf(Configuration.SuperserverUrlTemplate, "http", "%s", "8000", "%s")

	case TEST:
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		Configuration.SuperserverUrlTemplate = fmt.Sprintf(Configuration.SuperserverUrlTemplate, "http", "%s", "8000", "%s")

	case PRO:
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
		slog.SetDefault(logger)
		Configuration.SuperserverUrlTemplate = fmt.Sprintf(Configuration.SuperserverUrlTemplate, "http", "%s", "80", "%s")
	}

	slog.Info(fmt.Sprintf("Running in %s environment", environment))
}

func CloseConfig() {
	if fileDescriptor != nil {
		fileDescriptor.Close()
	}

}
