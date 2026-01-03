package context

type UserRights string

const ADMIN UserRights = "admin"
const RESERVATOR UserRights = "reservator"
const VISITOR UserRights = "visitor"
const SPECTATOR UserRights = "spectator"

type Stage int

const CHOOSING Stage = 0
const VOTING Stage = 1
const REMINDING Stage = 2
const APPROVING Stage = 3
const RESERVATING Stage = 4
const REVIEWING Stage = 5

type User struct {
	id       int
	username string
	rights   UserRights
}

func NewUser(id int, username string, rights UserRights) User {
	return User{id, username, rights}
}

type Context struct {
	user  User
	stage Stage
}

func NewContext(user User, stage Stage) *Context {
	return &Context{user, stage}
}
