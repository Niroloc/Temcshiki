package tasks

import (
	"github.com/Niroloc/Temcshiki/v2/src/context"
	"github.com/Niroloc/Temcshiki/v2/src/db"
)

type ChooseReminder struct {
	stages map[db.Stage]struct{}
}

func NewChooseReminder() *ChooseReminder {
	return &ChooseReminder{map[db.Stage]struct{}{db.CHOOSING: {}}}
}

func (this *ChooseReminder) Stages() map[db.Stage]struct{} {
	return this.stages
}

func (this *ChooseReminder) Apply(context *context.Context) {
	bot := context.GetBot()
	for userTgId, user := range context.GetUsers() {
		if user.Rights == db.ADMIN || user.Rights == db.RESERVATOR || user.Rights == db.VISITOR {
			bot.SendMessage(userTgId,
				"Доброе напоминание о том, что сейчас самое время добавить вариантов для будущих мероприятий!")
		}
		if user.Rights == db.ADMIN {
			bot.SendMessage(userTgId, "А лично тебе было бы неплохо обозначить даты для будущих выборов!")
		}
	}
}
