package db

import (
	"database/sql"
	"fmt"

	"github.com/Niroloc/Temcshiki/v2/src/context"
)

type Db struct {
	connection *sql.DB
}

func GetDb(file string) *Db {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
	return &Db{db}
}

func (this *Db) GetContext(userId int) *context.Context {
	res, err := this.connection.Exec(
		fmt.Sprint("SELECT * FROM users where userId = %d",
			userId))
	if err != nil {
		panic(err)
	}
	line, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	fmt.Printf("line: %v\n", line)
	user := context.NewUser(0, "Konfeen Ott", context.ADMIN)
	stage := context.CHOOSING
	return context.NewContext(user, stage)
}
