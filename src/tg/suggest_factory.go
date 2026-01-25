package tg

import (
	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type SuggestFactory struct {
	alias  data.Alias
	logger *logger.Logger
	args   SuggestFactoryArgs
}

type SuggestFactoryArgs struct {
	RestName     string
	MapUrlId     int
	ClosestMetro string
}
