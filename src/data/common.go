package data

import "time"

type CommonData struct {
	stage        Stage
	nextTaskDate string
	users        map[int]*User
	events       map[int]Event
	rests        map[int]Rest
	reviews      map[int]Review
	dates        map[int]Date
	votes        map[int]Vote
}

func MakeCommonData(dbWrapper *Db, tgIdToUser map[int]*User) *CommonData {
	result := &CommonData{
		stage:        dbWrapper.GetStage(),
		nextTaskDate: dbWrapper.GetNextTaskDate(),
		users:        map[int]*User{},
		events:       map[int]Event{},
		rests:        map[int]Rest{},
		reviews:      map[int]Review{},
		dates:        map[int]Date{},
		votes:        map[int]Vote{},
	}
	for _, user := range tgIdToUser {
		result.users[user.Id] = user
	}
	for _, event := range dbWrapper.GetEvents() {
		id := event.id
		visitDate := time.Now()
		if event.visitDate.Valid {
			visitDate, _ = time.Parse(time.DateOnly, event.visitDate.String)
		}
		restId := -1
		if event.restId.Valid {
			restId = int(event.restId.Int64)
		}
		result.events[id] = Event{id, visitDate, restId}
	}
	//ToDo: continue this shit
	return result
}
