package data

import "github.com/Niroloc/Temcshiki/v2/src/utils"

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
	result.events = utils.ListToMap(
		utils.ListMap(dbWrapper.GetEvents(), EventFromExported),
		func(e Event) int { return e.Id },
	)
	for _, rest := range dbWrapper.GetRests() {
		result.rests[rest.Id] = rest
	}
	for _, review := range dbWrapper.GetReviews() {
		result.reviews[review.Id] = review
	}
	for _, exportedDate := range dbWrapper.GetDates() {
		date := DateFromExported(exportedDate)
		result.dates[date.Id] = date
	}
	for _, exportedVote := range dbWrapper.GetVotes() {
		vote := VoteFromExported(exportedVote)
		result.votes[vote.Id] = vote
	}
	return result
}

func (this *CommonData) CreateNewEvent(e Event) {
	this.events[e.Id] = e
}
