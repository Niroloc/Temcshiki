package tasks

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Reservating struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewReservating() *Reservating {
	return &Reservating{stages: map[data.Stage]struct{}{data.RESERVATING: {}},
		logger: logger.GetLogger(reflect.TypeFor[Reservating]())}
}

func (this *Reservating) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *Reservating) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting reservating cron task")
	event, err := exportData.GetNewEvent()
	if err != nil {
		this.logger.Error(err.Error())
		return
	}
	acceptedVisitors := exportData.GetAcceptedVisitors(event)
	visitorsNumber := len(acceptedVisitors)
	if event.RestId == -1 {
		this.logger.Error("Rest id is not correct in new event")
		return
	}
	rest := exportData.GetRestById(event.RestId)
	if event.VisitDate == nil {
		this.logger.Error("Visit date is not set up in the new event")
		return
	}
	date := *event.VisitDate
	for _, user := range exportData.GetUsers() {
		if user.Rights == data.RESERVATOR || user.Rights == data.ADMIN {
			err = bot.SendMessage(
				user,
				fmt.Sprintf(
					"Приветствую, дорогой бронировальщик! Нам нужно забронировать ресторан %s (%s) на %d человек %s в 16:00",
					rest.RestName,
					rest.MapUrl,
					visitorsNumber+1,
					date.Format(time.DateOnly),
				),
			)
			if err != nil {
				this.logger.Error("Message is not delivered")
			}
		}
	}
}
