package data

import (
	"errors"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type Data struct {
	tgIdToUser map[int]User
	stage      Stage
	db         *Db
	logger     *logger.Logger
}

func CreateData(dbWrapper *Db) *Data {
	tgIdToUser := map[int]User{}
	stage := dbWrapper.GetStage()
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, UserRights(user.Rights)}
	}
	return &Data{
		tgIdToUser: tgIdToUser,
		stage:      stage,
		db:         dbWrapper,
		logger:     logger.GetLogger(reflect.TypeFor[Data]())}
}

func (this *Data) GetStage() Stage {
	return this.stage
}

func (this *Data) GetUser(tgId int) (User, error) {
	if user, exists := this.tgIdToUser[tgId]; exists {
		return user, nil
	}
	this.logger.Warn("Unknown user has been searched")
	return User{0, "", ""}, errors.New("Unknown user")
}

func (this *Data) GetUsers() map[int]User {
	return this.tgIdToUser
}

func (this *Data) NextStage() {
	prevStage := this.stage
	curStage := this.stage.Next()
	this.db.UpdateStage(prevStage, curStage)
}

func (this *Data) GetDb() *Db {
	return this.db
}

func (this *Data) GetUsersMap() map[int]User {
	return this.tgIdToUser
}

func (this *Data) GetRestsForVoting() []Rest {
	return this.db.GetQueuedRests()
}
