package main

import (
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/aquasecurity/table"
	"github.com/joho/godotenv"
)

type root struct{}

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn(".env is not loaded, be careful!")
	}
	testMain()
	// dbWrapper := data.GetDb(os.Getenv("DB_FILE"))
	// dbWrapper.InitDb(os.Getenv("FORWARD_MIGRATION"))
	// logger.Info("Db initialized")
	// logger.Info("Bot initialized")
	// data := data.CreateData(dbWrapper)
	// bot := tg.CreateBot(os.Getenv("TOKEN"))
	// logger.Info("Data created")
	// tasks := tasks.InitTasks(data, bot)
	// logger.Info("Tasks initialized")
	// go tasks.Loop()
	// logger.Info("Tasks loop started")
	// err = bot.InfinitePolling(data)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}

func testMain() {
	t := table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("ID", "Fruit", "Stock")
	t.AddRow("1", "Apple", "14")
	t.AddRow("2", "Banana", "88,041")
	t.AddRow("3", "Cherry", "342")
	t.AddRow("4", "Dragonfruit", "1")
	t.Render()
}
