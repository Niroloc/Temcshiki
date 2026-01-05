package context

import "github.com/Niroloc/Temcshiki/v2/src/db"

type UserRights string

const ADMIN UserRights = "admin"
const RESERVATOR UserRights = "reservator"
const VISITOR UserRights = "visitor"
const SPECTATOR UserRights = "spectator"

type User struct {
	id       int
	username string
	rights   UserRights
}

func NewUser(id int, username string, rights UserRights) User {
	return User{id, username, rights}
}

type Context struct {
	tgIdToUser map[int]User
	stage      db.Stage
	fresh      bool
	db         *db.Db
}

func CreateContext(db *db.Db) *Context {
	tgIdToUser := map[int]User{}
	stage := db.GetStage()
	for _, user := range db.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, UserRights(user.Rights)}
	}
	return &Context{tgIdToUser: tgIdToUser, stage: stage, fresh: true, db: db}
}
