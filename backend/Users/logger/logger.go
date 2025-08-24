package logger

import (
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	writerhook "github.com/sirupsen/logrus/hooks/writer"
)

var Log = logrus.New()
var currentLevel atomic.Value // stores logrus.Level

// Init reads LOG_LEVEL + LOG_FORMAT and sets outputs.
// Levels shown on terminal are controlled by SetLevel(...) or env.
func Init() {
	level := parseLevel(os.Getenv("LOG_LEVEL")) // debug|info|warn|error|panic|fatal|trace
	SetLevel(level)

	// format
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "text" {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	} else {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// send info/debug/trace to stdout, warn/error/fatal/panic to stderr
	Log.SetOutput(io.Discard) // prevent double printing
	Log.AddHook(&writerhook.Hook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel,
		},
	})
	Log.AddHook(&writerhook.Hook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel,
		},
	})
}

func parseLevel(s string) logrus.Level {
	l, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(s)))
	if err != nil {
		return logrus.InfoLevel
	}
	return l
}

// SetLevel can be called at runtime (from admin API).
func SetLevel(l logrus.Level) {
	Log.SetLevel(l)
	currentLevel.Store(l)
}

func GetLevel() logrus.Level {
	if v := currentLevel.Load(); v != nil {
		return v.(logrus.Level)
	}
	return logrus.InfoLevel
}

// helper
func With(fields logrus.Fields) *logrus.Entry { return Log.WithFields(fields) }
