package data

type UserRights string

const ADMIN UserRights = "admin"
const RESERVATOR UserRights = "reservator"
const VISITOR UserRights = "visitor"
const SPECTATOR UserRights = "spectator"

func RightsIsCorrect(r UserRights) bool {
	if r != ADMIN && r != RESERVATOR && r != VISITOR && r != SPECTATOR {
		return false
	}
	return true
}
