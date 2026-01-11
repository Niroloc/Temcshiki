package data

type Stage int

const (
	CHOOSING    Stage = 0
	VOTING      Stage = 1
	COUNTING    Stage = 2
	REMINDING   Stage = 3
	RESERVATING Stage = 4
	REVIEWING   Stage = 5
)

func (this *Stage) Next() Stage {
	(*this) = ((*this) + 1) % 6
	return (*this)
}
