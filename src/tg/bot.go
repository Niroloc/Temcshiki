package tg

import (
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type Bot struct {
	logger *logger.Logger
}

func InitBot() *Bot {
	return &Bot{
		logger: logger.GetLogger(reflect.TypeFor[Bot]()),
	}
}

func (this *Bot) SendMessage(tgId int, msg string) error {
	return nil
}
