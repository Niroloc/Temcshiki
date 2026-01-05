package db

import (
	"database/sql"
	"os"
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

func (this *Db) InitDb() {
	sql_file := os.Getenv("FORWARD_MIGRATION")
	query, err := os.ReadFile(sql_file)
	if err != nil {
		panic(err)
	}
	if _, err := this.connection.Exec(string(query)); err != nil {
		panic(err)
	}
	res, err := this.connection.Query("SELECT current_stage FROM stage")
	if err != nil {
		panic(err)
	}
	if !res.Next() {
		this.connection.Exec("INSERT INTO stage (current_stage) VALUES (0)")
	}
}
