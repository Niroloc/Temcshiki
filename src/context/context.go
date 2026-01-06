package context

import "github.com/Niroloc/Temcshiki/v2/src/db"

type Context struct {
	tgIdToUser map[int]User
	stage      db.Stage
	fresh      bool
	db         *db.Db
}

func CreateContext(dbWrapper *db.Db) *Context {
	tgIdToUser := map[int]User{}
	stage := dbWrapper.GetStage()
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, db.UserRights(user.Rights)}
	}
	return &Context{tgIdToUser: tgIdToUser, stage: stage, fresh: true, db: dbWrapper}
}
