package logger

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type LogLevel int8

const DEBUG LogLevel = 0
const INFO LogLevel = 1
const WARN LogLevel = 2
const ERROR LogLevel = 3

type Logger struct {
	logger *log.Logger
	level  LogLevel
	source string
}

func getLogLevel() LogLevel {
	str := os.Getenv("LOGLEVEL")
	if str == "" {
		return WARN
	}
	intLevel, err := strconv.Atoi(str)
	if err != nil {
		return WARN
	}
	if intLevel < 0 || intLevel > 3 {
		return DEBUG
	}
	level := LogLevel(intLevel)
	return level
}

func GetLogger(source reflect.Type) *Logger {
	return &Logger{log.Default(), getLogLevel(), source.Name()}
}

func (this *Logger) Debug(msg string) {
	if this.level > DEBUG {
		return
	}
	this.logger.Output(6, strings.Join([]string{"[DEBUG]", this.source, ":", msg}, " "))
}

func (this *Logger) Info(msg string) {
	if this.level > INFO {
		return
	}
	this.logger.Output(6, strings.Join([]string{"[INFO]", this.source, ":", msg}, " "))
}

func (this *Logger) Warn(msg string) {
	if this.level > WARN {
		return
	}
	this.logger.Output(6, strings.Join([]string{"[WARN]", this.source, ":", msg}, " "))
}

func (this *Logger) Error(msg string) {
	if this.level > ERROR {
		return
	}
	this.logger.Output(6, strings.Join([]string{"[ERROR]", this.source, ":", msg}, " "))
}
