package data

import (
	"fmt"
	"time"
)

type Date struct {
	Id        int
	Candidate time.Time
	EventId   int
}

func DateFromExported(exportedDate ExportedDate) Date {
	candidate, _ := time.Parse(time.DateOnly, exportedDate.candidate)
	return Date{exportedDate.id, candidate, exportedDate.eventId}
}

func (this Date) GetDescription(i int) string {
	return fmt.Sprintf(
		"%d) %s\n",
		i+1,
		this.Candidate.Format(time.DateOnly),
	)
}

func (this Date) GetButtonTitle() string {
	return this.Candidate.Format(time.DateOnly)
}

func (this Date) GetCallbackData(eventId int) string {
	return fmt.Sprintf("%s_%d_%d", DATE, eventId, this.Id)
}
