package context

import "github.com/Niroloc/Temcshiki/v2/src/db"

type User struct {
	Id       int
	Username string
	Rights   db.UserRights
}

func NewUser(id int, username string, rights db.UserRights) User {
	return User{id, username, rights}
}
