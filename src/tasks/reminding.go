package tasks

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Reminding struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewReminding() *Reminding {
	return &Reminding{stages: map[data.Stage]struct{}{data.REMINDING: {}},
		logger: logger.GetLogger(reflect.TypeFor[Reminding]())}
}

func (this *Reminding) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *Reminding) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting reminding cron task")
	event, err := exportData.GetNewEvent()
	if err != nil {
		this.logger.Error("What?!")
		return
	}
	objs := []data.VotingObject{event}
	rest := exportData.GetRestById(event.RestId)
	for userTgId, user := range exportData.GetUsers() {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			bot.SendMessageWithVoting(
				userTgId,
				fmt.Sprintf(
					"Прошу прожать кнопочку в знак подтверждения посещения мероприятия %s в 16:00 в %s\nВстреча в метро %s в 15:45",
					event.VisitDate.Format(time.DateOnly),
					rest.RestName,
					rest.ClosestMetro,
				),
				event.Id,
				objs,
			)
		}
	}
}
