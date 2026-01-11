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
	if curStage == CHOOSING {
		nextDate = getNextWeekday(curDate, time.Monday)
	} else if curStage == VOTING {
		nextDate = getNextWeekday(curDate, time.Sunday)
	} else if curStage == COUNTING {
		nextDate = getNextWeekday(curDate, time.Sunday)
	} else if curStage == REMINDING {
		eventDate := *this.commonData.GetNextEvent().visitDate
		nextDate = getPrevWeekDay(eventDate, time.Monday)
	} else if curStage == RESERVATING {
		nextDate = getNextWeekday(curDate, time.Tuesday)
	} else if curStage == REVIEWING {
		nextDate = *this.commonData.GetNextEvent().visitDate
	} else {
		nextDate = time.Now().AddDate(0, 0, 1)
	}
	this.db.UpdateNextTaskDate(nextDate)
}

func (this *Data) GetRestsForVoting() []Rest {
	visitedRestIds := map[int]struct{}{}
	for _, e := range this.commonData.events {
		if e.restId > -1 {
			visitedRestIds[e.restId] = struct{}{}
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
			return e.restId == -1
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
	month := time.Now().Month()
	year := time.Now().Year()
	year += int(month) / 12
	month = month%12 + 1
	return utils.FilterValues(
		this.commonData.dates,
		func(d Date) bool {
			return d.Candidate.Month() == month && d.Candidate.Year() == year
		},
	)
}

func (this *Data) GetNextTaskDate() string {
	return this.db.GetNextTaskDate()
}
