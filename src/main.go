package main

import (
	"os"
	"strconv"

	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

func getLogLevel() logger.LogLevel {
	str := os.Getenv("LOGLEVEL")
	if str == "" {
		return logger.WARN
	}
	intLevel, err := strconv.Atoi(str)
	if err != nil {
		return logger.WARN
	}
	if intLevel < 0 || intLevel > 3 {
		return logger.DEBUG
	}
	level := logger.LogLevel(intLevel)
	return level
}

func main() {
	logger := logger.GetLogger(getLogLevel())
	dbWrapper := db.GetDb("data/db.db")
	dbWrapper.InitDb()
	logger.Debug("Db initialized successfully!")
}
