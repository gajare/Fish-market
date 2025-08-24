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

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Log.Fatal("DATABASE_URL not set")
	}

	sqlLevel := glogger.Warn
	switch os.Getenv("LOG_SQL_LEVEL") { // "silent","error","warn","info"
	case "silent":
		sqlLevel = glogger.Silent
	case "error":
		sqlLevel = glogger.Error
	case "info":
		sqlLevel = glogger.Info
	}
	slowMs := 200
	if v := os.Getenv("LOG_SQL_SLOW_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			slowMs = n
		}
	}

	// pipe GORM logs into logrus via std log writer
	std := log.New(logger.Log.Writer(), "", 0)
	gormLogger := glogger.New(
		std,
		glogger.Config{
			SlowThreshold:             time.Duration(slowMs) * time.Millisecond,
			LogLevel:                  sqlLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		logger.Log.Fatalf("gorm open: %v", err)
	}

	if err := gdb.AutoMigrate(&models.User{}); err != nil {
		logger.Log.Fatalf("automigrate: %v", err)
	}
	DB = gdb
	logger.With(map[string]any{"dsn": dsn}).Info("db_connected")
}
