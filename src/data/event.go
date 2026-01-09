package data

import "time"

type Event struct {
	Id        int
	visitDate time.Time
	restId    int
}

func GetEmptyEvent() Event {
	return Event{-1, time.Now(), -1}
}
