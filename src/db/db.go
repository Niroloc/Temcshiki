package db

import (
	"github.com/Niroloc/Temcshiki/v2/src/context"
)

type ExportedUser struct {
	id int
	tgid int
	username string
	rights string
}

func GetContext() context.Context {
	
}