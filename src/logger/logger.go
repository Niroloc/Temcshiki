package logger

import (
	"log"
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
}

func GetLogger(level LogLevel) *Logger {
	return &Logger{log.Default(), level}
}

func (this *Logger) Debug(msg string) {
	if this.level > DEBUG {
		return
	}
	this.logger.Output(2, strings.Join([]string{"[DEBUG]", msg}, " "))
}

func (this *Logger) Info(msg string) {
	if this.level > INFO {
		return
	}
	this.logger.Output(2, strings.Join([]string{"[INFO]", msg}, " "))
}

func (this *Logger) Warn(msg string) {
	if this.level > WARN {
		return
	}
	this.logger.Output(2, strings.Join([]string{"[WARN]", msg}, " "))
}

func (this *Logger) Error(msg string) {
	if this.level > ERROR {
		return
	}
	this.logger.Output(2, strings.Join([]string{"[ERROR]", msg}, " "))
}
