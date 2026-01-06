package tasks

import "github.com/Niroloc/Temcshiki/v2/src/db"

type Task interface {
	Stages() map[db.Stage]struct{}
	Apply()
}
