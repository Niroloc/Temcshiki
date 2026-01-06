package context

import "github.com/Niroloc/Temcshiki/v2/src/db"

type User struct {
	id       int
	username string
	rights   db.UserRights
}

func NewUser(id int, username string, rights db.UserRights) User {
	return User{id, username, rights}
}
