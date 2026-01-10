package tasks

import (
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type ChooseReminder struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewChooseReminder() *ChooseReminder {
	return &ChooseReminder{stages: map[data.Stage]struct{}{data.CHOOSING: {}},
		logger: logger.GetLogger(reflect.TypeFor[ChooseReminder]())}
}

func (this *ChooseReminder) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *ChooseReminder) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting choose_reminder cron task")
	for userTgId, user := range exportData.GetUsers() {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			err := bot.SendMessage(userTgId,
				"Доброе напоминание о том, что сейчас самое время добавить вариантов для будущих мероприятий!")
			if err != nil {
				this.logger.Warn("Seems like mesage to all is not delivered")
			}
		}
		if user.Rights == data.ADMIN {
			err := bot.SendMessage(userTgId,
				"А лично тебе было бы неплохо обозначить даты для будущих выборов!")
			if err != nil {
				this.logger.Warn("Message to admin is not delivered")
			}
		}
	}
}
