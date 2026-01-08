package data

type User struct {
	Id       int
	Username string
	Rights   UserRights
}

func NewUser(id int, username string, rights UserRights) User {
	return User{id, username, rights}
}
