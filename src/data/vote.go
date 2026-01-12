package data

type Vote struct {
	Id      int
	UserId  int
	EventId int
	RestId  int
	DateId  int
}

func VoteFromExported(exportedVote ExportedVote) Vote {
	vote := Vote{exportedVote.id, exportedVote.userId, exportedVote.eventId, -1, -1}
	if exportedVote.restId.Valid {
		vote.RestId = int(exportedVote.restId.Int64)
	}
	if exportedVote.dateId.Valid {
		vote.DateId = int(exportedVote.dateId.Int64)
	}
	return vote
}
