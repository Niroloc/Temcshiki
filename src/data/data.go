package data

import (
	"errors"
	"reflect"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/utils"
)

func getNextWeekday(today time.Time, weekday time.Weekday) time.Time {
	res := today.AddDate(0, 0, 1)
	for res.Weekday() != weekday {
		res = res.AddDate(0, 0, 1)
	}
	return res
}

func getPrevWeekDay(today time.Time, weekday time.Weekday) time.Time {
	res := today.AddDate(0, 0, -1)
	for res.Weekday() != weekday {
		res = res.AddDate(0, 0, -1)
	}
	return res
}

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

func (this *Data) GetCommonData() *CommonData {
	return this.commonData
}

func (this *Data) NextStage() {
	prevStage := this.commonData.stage
	curStage := this.commonData.stage.Next()
	this.db.UpdateStage(prevStage, curStage)
	curDate, _ := time.Parse(time.DateOnly, this.commonData.nextTaskDate)
	var nextDate time.Time
	switch curStage {
	case CHOOSING:
		nextDate = getNextWeekday(curDate, time.Monday)
	case COUNTING:
		nextDate = getNextWeekday(curDate, time.Sunday)
	case REMINDING:
		nextDate = getPrevWeekDay(*this.commonData.GetNextEvent().VisitDate, time.Monday)
	case RESERVATING:
		nextDate = getNextWeekday(curDate, time.Tuesday)
	case REVIEWING:
		nextDate = *this.commonData.GetNextEvent().VisitDate
	default:
		nextDate = time.Now().AddDate(0, 0, 1)
	}
	this.db.UpdateNextTaskDate(nextDate)
}

func (this *Data) GetRestsForVoting() []Rest {
	visitedRestIds := map[int]struct{}{}
	for _, e := range this.commonData.events {
		if e.RestId > -1 {
			visitedRestIds[e.RestId] = struct{}{}
		}
	}
	return utils.FilterValues(
		this.commonData.rests,
		func(r Rest) bool {
			_, exists := visitedRestIds[r.Id]
			return !exists
		},
	)
}

func (this *Data) GetNewEvent() (Event, error) {
	evs := utils.FilterValues(
		this.commonData.events,
		func(e Event) bool {
			return e.RestId == -1
		},
	)
	if len(evs) == 0 {
		return Event{}, errors.New("No new events")
	}
	return evs[0], nil
}

func (this *Data) CreateNewEvent() {
	e, err := this.db.CreateNewEvent()
	if err != nil {
		return
	}
	this.commonData.CreateNewEvent(e)
}

func (this *Data) GetDatesForVoting() []Date {
	ev, err := this.GetNewEvent()
	if err != nil {
		this.logger.Error("No dates for the new event!")
		return []Date{}
	}
	return utils.FilterValues(
		this.commonData.dates,
		func(d Date) bool {
			return ev.Id == d.EventId
		},
	)
}

func (this *Data) GetNextTaskDate() string {
	return this.db.GetNextTaskDate()
}

func (this *Data) GetWinnerDate(event Event) Date {
	votesByDateId := map[int]int{}
	for _, v := range utils.FilterValues(
		this.commonData.votes,
		func(v Vote) bool {
			return v.EventId == event.Id && v.DateId > -1
		},
	) {
		votesByDateId[v.DateId]++
	}
	winnerId := -1
	for id, votes := range votesByDateId {
		if votes > votesByDateId[winnerId] {
			winnerId = id
		}
	}
	return this.commonData.dates[winnerId]
}

func (this *Data) GetWinnerRest(event Event) Rest {
	votesByRestId := map[int]int{}
	for _, v := range utils.FilterValues(
		this.commonData.votes,
		func(v Vote) bool {
			return v.EventId == event.Id && v.RestId > -1
		},
	) {
		votesByRestId[v.RestId]++
	}
	winnerId := -1
	for id, votes := range votesByRestId {
		if votes > votesByRestId[winnerId] {
			winnerId = id
		}
	}
	return this.commonData.rests[winnerId]
}
