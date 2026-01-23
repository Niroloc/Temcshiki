package tasks

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type ReviewOption struct {
	val      int
	category data.ReviewCategory
}

func (this ReviewOption) GetDescription(i int) string {
	return ""
}

func (this ReviewOption) GetButtonTitle() string {
	return strconv.Itoa(this.val)
}

func (this ReviewOption) GetCallbackData(eventId int) string {
	return fmt.Sprintf("%s_%d_%d_%s", data.REVIEW, eventId, this.val, string(this.category))
}

type Reviewing struct {
	stages map[data.Stage]struct{}
	logger *logger.Logger
}

func NewReviewing() *Reviewing {
	return &Reviewing{stages: map[data.Stage]struct{}{data.REVIEWING: {}},
		logger: logger.GetLogger(reflect.TypeFor[Reviewing]())}
}

func (this *Reviewing) Stages() map[data.Stage]struct{} {
	return this.stages
}

func (this *Reviewing) Apply(bot *tg.Bot, exportData *data.Data) {
	this.logger.Info("Starting reviewing cron task")
	event, err := exportData.GetNewEvent()
	if err != nil {
		this.logger.Error(err.Error())
		return
	}
	interiors := []data.VotingObject{}
	services := []data.VotingObject{}
	foods := []data.VotingObject{}
	prices := []data.VotingObject{}
	for val := 1; val < 11; val++ {
		interiors = append(interiors, ReviewOption{val, data.INTERIOR})
		services = append(services, ReviewOption{val, data.SERVICE})
		foods = append(foods, ReviewOption{val, data.FOOD})
		prices = append(prices, ReviewOption{val, data.PRICES})
	}
	acceptedVisitors := exportData.GetAcceptedVisitors(event)
	for _, user := range acceptedVisitors {
		if user.Rights == data.ADMIN || user.Rights == data.RESERVATOR || user.Rights == data.VISITOR {
			err := bot.SendMessageWithVoting(&user, "Настаёт время голосовать! ГОЛОСОВАНИЕ!\nИнтерьер:", event.Id, interiors)
			if err != nil {
				this.logger.Error("Message has not been delievered")
			}
			err = bot.SendMessageWithVoting(&user, "Обслуживание:", event.Id, services)
			if err != nil {
				this.logger.Error("Message has not been delievered")
			}
			err = bot.SendMessageWithVoting(&user, "Цены:", event.Id, prices)
			if err != nil {
				this.logger.Error("Message has not been delievered")
			}
			err = bot.SendMessageWithVoting(&user, "Еда:", event.Id, foods)
			if err != nil {
				this.logger.Error("Message has not been delievered")
			}
		}
	}
}
