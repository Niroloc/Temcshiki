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

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	dbWrapper := data.GetDb(os.Getenv("DB_FILE"))
	dbWrapper.InitDb(os.Getenv("FORWARD_MIGRATION"))
	logger.Info("Db initialized")
	logger.Info("Bot initialized")
	data := data.CreateData(dbWrapper)
	bot := tg.CreateBot(os.Getenv("TOKEN"), data)
	logger.Info("Data created")
	tasks := tasks.InitTasks(data, bot)
	logger.Info("Tasks initialized")
	go tasks.Loop()
	logger.Info("Tasks loop started")
	err := bot.InfinitePolling(data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
