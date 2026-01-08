package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type root struct{}

func initTasks() []tasks.Task {
	return []tasks.Task{tasks.NewChooseReminder()}
}

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	dbWrapper := data.GetDb(os.Getenv("DB_FILE"))
	dbWrapper.InitDb(os.Getenv("FORWARD_MIGRATION"))
	logger.Info("Db initialized")
	bot := tg.CreateBot(os.Getenv("TOKEN"))
	logger.Info("Bot initialized")
	data := data.CreateData(dbWrapper)
	logger.Info("Data created")
	// tasks := initTasks()
	logger.Info("Tasks initialized")
	err := bot.InfinitePolling(data)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
