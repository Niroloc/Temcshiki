package data

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
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
	tgIdToUser    map[int]*User
	factoryRights *FactoryRights
	commonData    *CommonData
	db            *Db
	logger        *logger.Logger
}

func CreateData(dbWrapper *Db) *Data {
	tgIdToUser := map[int]*User{}
	for _, exportedUser := range dbWrapper.GetUsers() {
		user := UserFromExported(exportedUser)
		tgIdToUser[user.TgId] = user
	}
	return &Data{
		tgIdToUser:    tgIdToUser,
		factoryRights: CreateFactoryRights(),
		commonData:    MakeCommonData(dbWrapper, tgIdToUser),
		db:            dbWrapper,
		logger:        logger.GetLogger(reflect.TypeFor[Data]())}
}

func (this *Data) CheckFactoryRights(alias Alias, user *User) bool {
	return this.factoryRights.CheckUserRights(alias, user)
}

func (this *Data) GetStage() Stage {
	return this.commonData.stage
}

func (this *Data) GetUserByTg(tgId int) (*User, error) {
	if user, exists := this.tgIdToUser[tgId]; exists {
		return user, nil
	}
	this.logger.Warn("Unknown user has been found")
	return nil, errors.New("Unknown user")
}

func (this *Data) GetUserById(id int) (*User, error) {
	if user, exists := this.commonData.users[id]; exists {
		return user, nil
	}
	this.logger.Warn("Unknown user has been found")
	return nil, errors.New("Unknown user")
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
	return utils.FilterValuesToList(
		this.commonData.rests,
		func(r Rest) bool {
			_, exists := visitedRestIds[r.Id]
			return !exists
		},
	)
}

func (this *Data) GetNewEvent() (Event, error) {
	evs := utils.FilterValuesToList(
		this.commonData.events,
		func(e Event) bool {
			return !e.IsVisited
		},
	)
	if len(evs) == 0 {
		return Event{}, errors.New("No new events")
	}
	if len(evs) > 1 {
		this.logger.Warn("There are more than one unvisited event")
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
	return utils.FilterValuesToList(
		this.commonData.dates,
		func(d Date) bool {
			return ev.Id == d.EventId
		},
	)
}

func (this *Data) GetNextTaskDate() string {
	return this.commonData.nextTaskDate
}

func (this *Data) GetWinnerDate(event Event) Date {
	votesByDateId := map[int]int{}
	for _, v := range utils.FilterValuesToList(
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
	date := this.commonData.dates[winnerId]
	{
		dateForCopy := date.Candidate
		event.VisitDate = &dateForCopy
	}
	this.commonData.events[event.Id] = event
	this.db.UpdateEvent(event.Id, event)
	return date
}

func (this *Data) GetWinnerRest(event Event) Rest {
	votesByRestId := map[int]int{}
	for _, v := range utils.FilterValuesToList(
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
	rest := this.commonData.rests[winnerId]
	event.RestId = rest.Id
	this.commonData.events[event.Id] = event
	this.db.UpdateEvent(event.Id, event)
	return rest
}

func (this *Data) GetRestById(id int) Rest {
	return this.commonData.rests[id]
}

func (this *Data) GetAcceptedVisitors(event Event) []User {
	return utils.Map(
		utils.FilterValuesToList(
			this.commonData.votes,
			func(v Vote) bool {
				return v.RestId == -1 && v.DateId == -1 && v.EventId == event.Id
			},
		),
		func(v Vote) User {
			return *this.commonData.users[v.UserId]
		},
	)
}

func (this *Data) CreateNewUser(tgid int, username string, rights UserRights) *User {
	user := this.db.CreateNewUser(tgid, username, rights)
	if user == nil {
		return nil
	}
	this.commonData.users[user.Id] = user
	this.tgIdToUser[user.TgId] = user
	return user
}

func (this *Data) UpdateUser(user *User, newUsername string, newRights UserRights) {
	cur := this.commonData.users[user.Id]
	cur.Rights = newRights
	cur.Username = newUsername
	this.db.UpdateUser(user.Id, cur)
}

func (this *Data) GetRating() (Table, error) {
	table := CreateNewTable([]string{"#", "Название", "Дата посещения", "Интерьер", "Обслуживание", "Еда", "Цены", "Рейтинг"})
	reviews := utils.ValuesToList(this.commonData.reviews)
	restToReviews := utils.GroupByMapReduce(
		reviews,
		func(r Review) Rest { return this.commonData.rests[r.RestorauntId] },
		func(r Review) []Review { return []Review{r} },
		func(a []Review, b []Review) []Review { return append(a, b...) },
	)
	type Stat struct {
		sum    int
		number int
	}
	restToCategoryStats := map[Rest]map[ReviewCategory]Stat{}
	for rest, reviews := range restToReviews {
		restToCategoryStats[rest] = utils.GroupByMapReduce(
			reviews,
			func(r Review) ReviewCategory { return r.Category },
			func(r Review) Stat { return Stat{r.Rate, 1} },
			func(a, b Stat) Stat { return Stat{a.sum + b.sum, a.number + b.number} },
		)
	}
	restToOverall := map[Rest]Stat{}
	for rest, stats := range restToCategoryStats {
		restToOverall[rest] = *utils.Reduce(
			utils.ValuesToList(stats),
			func(a, b Stat) Stat { return Stat{a.sum + b.sum, a.number + b.number} },
		)
	}
	type Line struct {
		rest          Rest
		categoryStats map[ReviewCategory]Stat
		overallStat   Stat
	}
	lines := []Line{}
	for rest := range restToCategoryStats {
		lines = append(lines, Line{rest, restToCategoryStats[rest], restToOverall[rest]})
	}
	slices.SortFunc(lines, func(x, y Line) int {
		a := float64(x.overallStat.sum) / float64(x.overallStat.number)
		b := float64(y.overallStat.sum) / float64(y.overallStat.number)
		if a > b {
			return -1
		} else if a == b {
			return 0
		} else {
			return 1
		}
	})
	for i, l := range lines {
		eventO := utils.FilterValuesToList(this.commonData.events, func(e Event) bool { return e.RestId == l.rest.Id })
		if len(eventO) != 1 {
			return table, fmt.Errorf("Some issue during searching the event")
		}
		table.AddLine([]string{
			fmt.Sprintf("%d", i+1),
			l.rest.RestName,
			eventO[0].VisitDate.Format(time.DateOnly),
			fmt.Sprintf("%.2f", float64(l.categoryStats[INTERIOR].sum)/float64(l.categoryStats[INTERIOR].number)),
			fmt.Sprintf("%.2f", float64(l.categoryStats[SERVICE].sum)/float64(l.categoryStats[SERVICE].number)),
			fmt.Sprintf("%.2f", float64(l.categoryStats[FOOD].sum)/float64(l.categoryStats[FOOD].number)),
			fmt.Sprintf("%.2fv", float64(l.categoryStats[PRICES].sum)/float64(l.categoryStats[PRICES].number)),
			fmt.Sprintf("%.2f", float64(l.overallStat.sum)/float64(l.overallStat.sum)),
		})
	}
	return table, nil
}
