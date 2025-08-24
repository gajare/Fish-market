package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	// Level: debug|info|warn|error (default: info)
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	lvl, err := logrus.ParseLevel(levelStr)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Log.SetLevel(lvl)

	// Format: json|text (default: json)
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "text" {
		Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	} else {
		Log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}
	// Where to write (default: stderr)
	Log.SetOutput(os.Stdout)
}

// Convenience: attach common fields
func With(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
