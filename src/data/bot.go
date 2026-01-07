package data

import (
	"context"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot    *telego.Bot
	logger *logger.Logger
}

func CreateBot(botToken string) *Bot {
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

func (this *Bot) SendMessage(tgId int, msg string) error {
	ctx := context.Background()
	chatId := tu.ID(int64(tgId))
	_, err := this.bot.SendMessage(ctx, tu.Message(chatId, msg))
	return err
}
