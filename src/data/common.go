package data

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
	for _, exportedEvent := range dbWrapper.GetEvents() {
		event := EventFromExported(exportedEvent)
		result.events[event.Id] = event
	}
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
