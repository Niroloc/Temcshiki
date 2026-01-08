package data

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"

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

func (this *Db) InitDb(migration_file string) {
	this.logger.Debug(fmt.Sprintf("Forward sql script file: %s", os.Getenv("FORWARD_MIGRATION")))
	query, err := os.ReadFile(migration_file)
	if err != nil {
		this.logger.Error("No initial sql script")
		panic(err)
	}
	if _, err := this.connection.Exec(string(query)); err != nil {
		this.logger.Error("Error occured while execute initial script")
		panic(err)
	}
	res, err := this.connection.Query("SELECT current_stage FROM stage")
	if err != nil {
		this.logger.Error("Stage check failed")
		panic(err)
	}
	if !res.Next() {
		this.logger.Warn("No stage in db, inserting stage CHOOSING")
		this.connection.Exec("INSERT INTO stage (current_stage) VALUES (0)")
	}
	res, err = this.connection.Query("SELECT * FROM users WHERE rights = \"admin\"")
	if err != nil {
		this.logger.Error("Error while checking admin user")
		panic(err)
	}
	if res.Next() {
		return
	}
	this.logger.Warn("No admin user, default admin set")
	defaultAdminId, err := strconv.Atoi(os.Getenv("DEFAULT_ADMIN_ID"))
	if err != nil {
		this.logger.Error("No valid default admin id in envs")
		panic(err)
	}
	defaultAdminName := os.Getenv("DEFAULT_ADMIN_NAME")
	if _, err := this.connection.Exec(
		fmt.Sprintf("INSERT INTO users (tgid, username, rights) VALUES (%d, \"%s\", \"%s\")",
			defaultAdminId, defaultAdminName, ADMIN)); err != nil {
		this.logger.Error("Error while inserting default admin")
		panic(err)
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
