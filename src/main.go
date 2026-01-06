package main

import (
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/context"
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type root struct{}

func initTasks() []context.Task {
	return []context.Task{tasks.NewChooseReminder()}
}

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	dbWrapper := db.GetDb(os.Getenv("DB_FILE"))
	dbWrapper.InitDb()
	logger.Info("Db initialized")
	bot := tg.InitBot()
	logger.Info("Bot initialized")
	context := context.CreateContext(bot, dbWrapper, initTasks())
	logger.Info("Context created")
	context.NextStage()
}
