package data

import (
	"errors"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/utils"
)

type Data struct {
	tgIdToUser map[int]*User
	commonData *CommonData
	db         *Db
	logger     *logger.Logger
}

func CreateData(dbWrapper *Db) *Data {
	tgIdToUser := map[int]*User{}
	for _, exportedUser := range dbWrapper.GetUsers() {
		user := UserFromExported(exportedUser)
		tgIdToUser[user.TgId] = user
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
	visitedRestIds := map[int]struct{}{}
	for _, e := range this.commonData.events {
		if e.restId > -1 {
			visitedRestIds[e.restId] = struct{}{}
		}
	}
	return utils.FilterMapValues(
		this.commonData.rests,
		func(r Rest) bool {
			_, exists := visitedRestIds[r.Id]
			return !exists
		},
	)
}

func (this *Data) GetNewEvent() (Event, error) {
	evs := utils.FilterMapValues(
		this.commonData.events,
		func(e Event) bool {
			return e.restId == -1
		},
	)
	if len(evs) == 0 {
		return Event{}, errors.New("No new events")
	}
	return evs[0], nil
}

func (this *Data) GetDatesForVoting() []Date {
	return []Date{} //ToDo: finish that
}
