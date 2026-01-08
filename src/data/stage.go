package data

type Stage int

const CHOOSING Stage = 0
const VOTING Stage = 1
const REMINDING Stage = 2
const APPROVING Stage = 3
const RESERVATING Stage = 4
const REVIEWING Stage = 5

func (this *Stage) Next() {
	(*this) = ((*this) + 1) % 6
}
