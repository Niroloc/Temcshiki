package tasks

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Counting struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewCounting() *StartVoting {
	return &StartVoting{stages: map[data.Stage]struct{}{data.COUNTING: {}},
		logger: logger.GetLogger(reflect.TypeFor[Counting]())}
}

func (this *Counting) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *Counting) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting counting cron task")
	event, err := exportData.GetNewEvent()
	if err != nil {
		this.logger.Error(err.Error())
		return
	}
	winnerDate := exportData.GetWinnerDate(event)
	winnerRest := exportData.GetWinnerRest(event)
	for _, user := range exportData.GetUsers() {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			bot.SendMessage(
				user,
				fmt.Sprintf(
					"Голосование завершилось! был выбран ресторан: %s, и дата: %s",
					winnerRest.RestName, winnerDate.Candidate.Format(time.DateOnly),
				),
			)
			bot.SendMessage(
				user,
				fmt.Sprintf(
					"Подробная информация по встрече:\nРесторан: %s\nСсылка: %s\nМетро для встречи в 15:45 : %s",
					winnerRest.RestName,
					winnerRest.MapUrl,
					winnerRest.ClosestMetro,
				),
			)
		}
		if user.Rights == data.SPECTATOR {
			bot.SendMessage(
				user,
				fmt.Sprintf(
					"Ребята идут в %s, %s",
					winnerRest.RestName,
					winnerDate.Candidate,
				),
			)
		}
	}
}
