package db

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/gajare/Fish-market/logger"
	"github.com/gajare/Fish-market/models"
)

var DB *gorm.DB
var gormLogger glogger.Interface

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Log.Fatal("DATABASE_URL not set")
	}

	gormLogger = makeGormLogger(envGormLevel(), EnvSlowMS())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		logger.Log.Fatalf("gorm open: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Log.Fatalf("automigrate: %v", err)
	}
	DB = db
	logger.With(map[string]any{"dsn": dsn}).Info("db_connected")
}

func makeGormLogger(level glogger.LogLevel, slowMs int) glogger.Interface {
	std := log.New(logger.Log.Writer(), "", 0) // pipe to logrus writer
	return glogger.New(std, glogger.Config{
		SlowThreshold:             time.Duration(slowMs) * time.Millisecond,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
}

func envGormLevel() glogger.LogLevel {
	switch os.Getenv("LOG_SQL_LEVEL") { // "silent","error","warn","info"
	case "silent":
		return glogger.Silent
	case "error":
		return glogger.Error
	case "info":
		return glogger.Info
	default:
		return glogger.Warn
	}
}

func EnvSlowMS() int {
	if v := os.Getenv("LOG_SQL_SLOW_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return 200
}

// SetSQLLogLevel lets you change SQL log level at runtime if you like.
func SetSQLLogLevel(level string, slowMs int) {
	gormLogger = makeGormLogger(stringToGormLevel(level), slowMs)
	DB.Config.Logger = gormLogger
	logger.With(map[string]any{"sql_level": level, "slow_ms": slowMs}).Info("gorm_logger_updated")
}

func stringToGormLevel(level string) glogger.LogLevel {
	switch level {
	case "silent":
		return glogger.Silent
	case "error":
		return glogger.Error
	case "info":
		return glogger.Info
	default:
		return glogger.Warn
	}
}
