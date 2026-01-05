package db

import (
	"database/sql"
	"os"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
	logger     *logger.Logger
}

type ExportedUser struct {
	Id       int
	TgId     int
	Username string
	Rights   string
}

func GetDb(file string) *Db {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
	return &Db{db, logger.GetLogger(reflect.TypeFor[Db]())}
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

func (this *Db) GetUsers() []ExportedUser {
	rows, err := this.connection.Query("SELECT * FROM users")
	if err != nil {
		this.logger.Warn("An error occured while getting users from DB")
		this.logger.Warn(err.Error())
		return []ExportedUser{}
	}
	var ans []ExportedUser
	for rows.Next() {
		user := ExportedUser{}
		err := rows.Scan(&user.Id, &user.TgId, &user.Username, &user.Rights)
		if err != nil {
			this.logger.Error("An error occured while parsing user")
			this.logger.Error(err.Error())
			panic(err)
		}
		ans = append(ans, user)
	}
	return ans
}

func (this *Db) GetStage() Stage {
	row, err := this.connection.Query("SELECT current_stage FROM stage")
	if err != nil || !row.Next() {
		this.logger.Warn("Cannot get stage from DB")
		this.logger.Warn(err.Error())
		return CHOOSING
	}
	var stage int
	row.Scan(&stage)
	return Stage(stage)
}
