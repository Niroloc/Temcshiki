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
	tempData         []any
}

func MakeUser(id int, tgId int, username string, rights UserRights) *User {
	return &User{id, tgId, username, rights, UserHistory{nil, false, []any{}}}
}

func UserFromExported(user ExportedUser) *User {
	return &User{
		Id:       user.Id,
		TgId:     user.TgId,
		Username: user.Username,
		Rights:   UserRights(user.Rights),
		History:  UserHistory{nil, false, []any{}},
	}
}
