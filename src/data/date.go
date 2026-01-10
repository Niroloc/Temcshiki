package data

import "time"

type Date struct {
	Id        int
	Candidate time.Time
	EventId   int
}

func DateFromExported(exportedDate ExportedDate) Date {
	candidate, _ := time.Parse(time.DateOnly, exportedDate.candidate)
	return Date{exportedDate.id, candidate, exportedDate.eventId}
}
