package context

import (
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
)

type Context struct {
	tgIdToUser map[int]User
	stage      db.Stage
	tasks      []tasks.Task
	db         *db.Db
}

func CreateContext(dbWrapper *db.Db, tasks []tasks.Task) *Context {
	tgIdToUser := map[int]User{}
	stage := dbWrapper.GetStage()
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, db.UserRights(user.Rights)}
	}
	return &Context{
		tgIdToUser: tgIdToUser,
		stage:      stage,
		tasks:      tasks,
		db:         dbWrapper}
}
