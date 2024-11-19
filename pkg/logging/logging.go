package logging

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Get the logging level from the environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default to info if not set
	}

	// Parse the log level
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'", logLevel)
		level = logrus.InfoLevel
	}

	log.SetLevel(level)
	return log
}

var Logger *logrus.Logger
