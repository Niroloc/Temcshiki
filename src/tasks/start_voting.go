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
	this.logger.Info("Starting start_voting cron task")
	event, err := exportData.GetNewEvent()
	if err != nil {
		this.logger.Error(err.Error())
		return
	}
	rests := []data.VotingObject{}
	for _, rest := range exportData.GetRestsForVoting() {
		rests = append(rests, rest)
	}
	dates := []data.VotingObject{}
	for _, date := range exportData.GetDatesForVoting() {
		dates = append(dates, date)
	}
	for _, user := range exportData.GetUsers() {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			err = bot.SendMessageWithVoting(
				user,
				"Стартуем наше головосание по ресторанам:",
				event.Id,
				rests,
			)
			if err != nil {
				this.logger.Warn("Seems like mesage with rests voting is not delivered")
			}
			err = bot.SendMessageWithVoting(
				user,
				"И по датам!",
				event.Id,
				dates,
			)
			if err != nil {
				this.logger.Warn("Seems like mesage with dates voting is not delivered")
			}
		}
	}
}
