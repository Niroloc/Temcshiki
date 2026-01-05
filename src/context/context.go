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
	tgidToUser map[int]User
	stage      db.Stage
	fresh      bool
	db         *db.Db
}

func NewContext(db *db.Db) *Context {
	return &Context{}
}
