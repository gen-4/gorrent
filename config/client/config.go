package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

var fileDescriptor *os.File

func getEnv() string {
	var err error
	environment := "dev"

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

func Config() {
	var environment string = getEnv()
	logFile := os.Getenv("LOG_FILE")
	f, err := tea.LogToFile(logFile, "debug")
	if err != nil {
		slog.Error("Error opening log file", "error", err.Error())
		os.Exit(1)
	}
	fileDescriptor = f

	logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
	slog.SetDefault(logger)

	switch environment {
	case "dev":
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)

	case "test":
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)

	case "pro":
		logger := slog.New(slog.NewJSONHandler(fileDescriptor, nil))
		slog.SetDefault(logger)
	}

	slog.Info(fmt.Sprintf("Running in %s environment", environment))
}

func CloseConfig() {
	if fileDescriptor != nil {
		fileDescriptor.Close()
	}

}
