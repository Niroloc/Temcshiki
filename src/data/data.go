package data

import (
	"errors"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type Data struct {
	tgIdToUser map[int]*User
	commonData *CommonData
	db         *Db
	logger     *logger.Logger
}

func CreateData(dbWrapper *Db) *Data {
	tgIdToUser := map[int]*User{}
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = MakeUser(user.Id, user.TgId, user.Username, UserRights(user.Rights))
	}
	return &Data{
		tgIdToUser: tgIdToUser,
		commonData: MakeCommonData(dbWrapper, tgIdToUser),
		db:         dbWrapper,
		logger:     logger.GetLogger(reflect.TypeFor[Data]())}
}

func (this *Data) GetStage() Stage {
	return this.commonData.stage
}

func (this *Data) GetUser(tgId int) (User, error) {
	if user, exists := this.tgIdToUser[tgId]; exists {
		return *user, nil
	}
	this.logger.Warn("Unknown user has been searched")
	return *MakeUser(-1, -1, "", SPECTATOR), errors.New("Unknown user")
}

func (this *Data) GetUsers() map[int]*User {
	return this.tgIdToUser
}

func (this *Data) NextStage() {
	prevStage := this.commonData.stage
	curStage := this.commonData.stage.Next()
	this.db.UpdateStage(prevStage, curStage)
}

func (this *Data) GetDb() *Db {
	return this.db
}

func (this *Data) GetRestsForVoting() []Rest {
	return this.db.GetQueuedRests()
}
