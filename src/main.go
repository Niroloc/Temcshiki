package main

import (
	"log"

	"github.com/Niroloc/Temcshiki/v2/src/db"
)

func main() {
	logger := log.Default()
	logger.Output(2, "Test log")
	dbWrapper := db.GetDb("data/db.db")
	dbWrapper.InitDb()
	logger.Output(2, "Db initialized successfully!")
}
