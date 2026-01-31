package data

import (
	"fmt"

	"github.com/Niroloc/Temcshiki/v2/src/utils"
)

type CommonData struct {
	stage        Stage
	nextTaskDate string
	users        map[int]*User
	events       map[int]Event
	rests        map[int]Rest
	reviews      map[int]Review
	dates        map[int]Date
	votes        map[int]Vote
	urls         map[int]Url
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
		urls:         map[int]Url{},
	}
	for _, user := range tgIdToUser {
		result.users[user.Id] = user
	}
	result.events = utils.ToMapMappingToKey(
		utils.Map(dbWrapper.GetEvents(), EventFromExported),
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
	for _, url := range dbWrapper.GetUrls() {
		result.urls[url.Id] = url
	}
	return result
}

func (this *CommonData) CreateNewEvent(e Event) {
	this.events[e.Id] = e
}

func (this *CommonData) GetNextEvent() Event {
	id := 0
	for i := range this.events {
		id = max(id, i)
	}
	return this.events[id]
}

func (this *CommonData) GetUrl(id int) (Url, error) {
	url, exists := this.urls[id]
	if exists {
		return url, nil
	}
	return url, fmt.Errorf("Url with id %d does not exist", id)
}
