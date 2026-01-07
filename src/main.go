package main

import (
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type root struct{}

func initTasks() []data.Task {
	return []data.Task{tasks.NewChooseReminder()}
}

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	dbWrapper := db.GetDb(os.Getenv("DB_FILE"))
	dbWrapper.InitDb(os.Getenv("FORWARD_MIGRATION"))
	logger.Info("Db initialized")
	bot := tg.InitBot(os.Getenv("TOKEN"))
	logger.Info("Bot initialized")
	data := data.CreateContext(bot, dbWrapper, initTasks())
	logger.Info("Data created")
	data.NextStage()
	logger.Info(bot.InfinitePolling().Error())
}
