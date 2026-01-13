package data

import "time"

type Event struct {
	Id        int
	VisitDate *time.Time
	RestId    int
	IsVisited bool
}

func EventFromExported(event ExportedEvent) Event {
	id := event.id
	var visitDate *time.Time
	if event.visitDate.Valid {
		visitDateVal, _ := time.Parse(time.DateOnly, event.visitDate.String)
		visitDate = &visitDateVal
	}
	restId := -1
	if event.restId.Valid {
		restId = int(event.restId.Int64)
	}
	return Event{id, visitDate, restId, event.isVisited == 1}
}
