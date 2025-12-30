package context

type User struct {
	userId int
	userName string
}

func NewUser(userId int, userName string) *User {
	return &User{
		userId: userId,
		userName: userName,
	}
}