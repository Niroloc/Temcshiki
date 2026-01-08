package tasks

import (
	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Task interface {
	Stages() map[data.Stage]struct{}
	Apply(*tg.Bot, *data.Data)
}
