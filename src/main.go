package main

import (
	"fmt"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/context"
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
)

type root struct{}

func main() {
	logger := logger.GetLogger(reflect.TypeFor[root]())
	dbWrapper := db.GetDb("data/db.db")
	dbWrapper.InitDb()
	logger.Info("Db initialized successfully!")
	context := context.CreateContext(dbWrapper, []tasks.Task{})
	logger.Info(fmt.Sprint(context))
}
