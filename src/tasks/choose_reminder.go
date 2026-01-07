package tasks

import (
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type ChooseReminder struct {
	stages map[db.Stage]struct{}
	logger *logger.Logger
}

func NewChooseReminder() *ChooseReminder {
	return &ChooseReminder{stages: map[db.Stage]struct{}{db.CHOOSING: {}},
		logger: logger.GetLogger(reflect.TypeFor[ChooseReminder]())}
}

func (this *ChooseReminder) Stages() map[db.Stage]struct{} {
	return this.stages
}

func (this *ChooseReminder) Apply(data *data.Data) {
	bot := data.GetBot()
	this.logger.Info("Starting cron task")
	for userTgId, user := range data.GetUsers() {
		if user.Rights == db.ADMIN || user.Rights == db.RESERVATOR || user.Rights == db.VISITOR {
			err := bot.SendMessage(userTgId,
				"Доброе напоминание о том, что сейчас самое время добавить вариантов для будущих мероприятий!")
			if err != nil {
				this.logger.Warn("Seems like mesage to all is not delivered")
			}
		}
		if user.Rights == db.ADMIN {
			err := bot.SendMessage(userTgId,
				"А лично тебе было бы неплохо обозначить даты для будущих выборов!")
			if err != nil {
				this.logger.Warn("Message to admin is not delivered")
			}
		}
	}
}
