package tg

import (
	"context"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
)

type Bot struct {
	bot    *telego.Bot
	logger *logger.Logger
}

func InitBot(botToken string) *Bot {
	logger := logger.GetLogger(reflect.TypeFor[Bot]())
	bot, err := telego.NewBot(botToken)
	if err != nil {
		logger.Error("Error while creating bot")
		panic(err)
	}
	return &Bot{
		bot:    bot,
		logger: logger,
	}
}

func (this *Bot) InfinitePolling() error {
	updates, err := this.bot.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		this.logger.Error("Cannot start polling")
		panic(err)
	}
}

func (this *Bot) SendMessage(tgId int, msg string) error {
	return nil
}
