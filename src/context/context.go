package context

import (
	"errors"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Context struct {
	tgIdToUser map[int]User
	stage      db.Stage
	tasks      []Task
	bot        *tg.Bot
	db         *db.Db
	logger     *logger.Logger
}

func CreateContext(bot *tg.Bot, dbWrapper *db.Db, tasks []Task) *Context {
	tgIdToUser := map[int]User{}
	stage := dbWrapper.GetStage()
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, db.UserRights(user.Rights)}
	}
	return &Context{
		tgIdToUser: tgIdToUser,
		stage:      stage,
		tasks:      tasks,
		bot:        bot,
		db:         dbWrapper,
		logger:     logger.GetLogger(reflect.TypeFor[Context]())}
}

func (this *Context) GetTasks() []Task {
	return this.tasks
}

func (this *Context) GetStage() db.Stage {
	return this.stage
}

func (this *Context) GetUser(tgId int) (User, error) {
	if user, exists := this.tgIdToUser[tgId]; exists {
		return user, nil
	}
	this.logger.Warn("Unknown user has been searched")
	return User{0, "", ""}, errors.New("Unknown user")
}

func (this *Context) GetUsers() map[int]User {
	return this.tgIdToUser
}

func (this *Context) GetBot() *tg.Bot {
	return this.bot
}

func (this *Context) NextStage() {
	this.stage.Next()
}
