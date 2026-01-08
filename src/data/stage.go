package data

type Stage int

const CHOOSING Stage = 0
const VOTING Stage = 1
const COUNTING Stage = 2
const APPROVING Stage = 3
const RESERVATING Stage = 4
const REVIEWING Stage = 5

func (this *Stage) Next() Stage {
	(*this) = ((*this) + 1) % 6
	return (*this)
}
