package tg

import (
	"context"
	"errors"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
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
	ctx := context.Background()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message
		chatId := tu.ID(msg.Chat.ID)
		this.bot.SendMessage(ctx, tu.Message(chatId, "Привет, я бот!"))
	}
	return errors.New("Polling closed")
}

func (this *Bot) SendMessage(tgId int, msg string) error {
	ctx := context.Background()
	chatId := tu.ID(int64(tgId))
	_, err := this.bot.SendMessage(ctx, tu.Message(chatId, msg))
	return err
}
