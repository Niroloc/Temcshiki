package data

type User struct {
	Id       int
	TgId     int
	Username string
	Rights   UserRights
	History  UserHistory
}

type UserHistory struct {
	LastCallbackData *string
	InputMode        bool
}

func MakeUser(id int, tgId int, username string, rights UserRights) *User {
	return &User{id, tgId, username, rights, UserHistory{nil, false}}
}
