package main

import (
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
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
	bot := data.CreateBot(os.Getenv("TOKEN"))
	logger.Info("Bot initialized")
	data := data.CreateData(dbWrapper, bot, initTasks())
	logger.Info("Data created")
	logger.Info(data.InfinitePolling().Error())
}
