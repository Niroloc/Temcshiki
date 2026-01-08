package tasks

import (
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type StartVoting struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewStartVoting() *StartVoting {
	return &StartVoting{stages: map[data.Stage]struct{}{data.VOTING: {}},
		logger: logger.GetLogger(reflect.TypeFor[StartVoting]())}
}

func (this *StartVoting) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *StartVoting) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting cron task")
	for userTgId, user := range exportData.GetUsers() {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			err := bot.SendMessageWithRests(
				userTgId,
				"Стартуем наше головосание по ресторанам:",
				exportData.GetRestsForVoting(),
			)
			if err != nil {
				this.logger.Warn("Seems like mesage to all is not delivered")
			}
			// err = bot.SendMessageWithDates(
			// 	userTgId,
			// 	"И по датам!",
			// 	exportData.GetDateForVoting(),
			// )
		}
	}
}
