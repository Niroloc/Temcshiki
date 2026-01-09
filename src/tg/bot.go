package tg

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/data"
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

func (this *Bot) SendMessageWithRests(tgId int, prefix string, rests []data.Rest) error {
	ctx := context.Background()
	chatId := tu.ID(int64(tgId))
	restsButtons := []telego.InlineKeyboardButton{}
	suffix := ""
	for i, rest := range rests {
		suffix += fmt.Sprintf(
			"%d) %s: %s, встретимся на станции %s\n",
			i+1,
			rest.RestName,
			rest.MapUrl,
			rest.ClosestMetro,
		)
		restsButtons = append(
			restsButtons,
			tu.InlineKeyboardButton(rest.RestName).WithCallbackData(fmt.Sprintf("rest_%d", rest.Id)),
		)
	}
	buttonRows := [][]telego.InlineKeyboardButton{}
	for i := 3; i < len(restsButtons); i += 3 {
		buttonRows = append(buttonRows, restsButtons[i-3:i])
	}
	inlineKeyboard := tu.InlineKeyboard(buttonRows...)
	if _, err := this.bot.SendMessage(
		ctx,
		tu.Message(chatId, strings.Join([]string{prefix, suffix}, "\n")).WithReplyMarkup(inlineKeyboard),
	); err != nil {
		return err
	}
	return nil
}

func (this *Bot) InfinitePolling(exportedData *data.Data) error {
	updates, err := this.bot.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		this.logger.Error("Cannot start polling")
		panic(err)
	}
	this.logger.Info("Starting polling")
	for update := range updates {
		if update.Message != nil {
			tgId := int(update.Message.Chat.ID)
			this.SendMessage(tgId, fmt.Sprint(exportedData.GetUsers()[tgId].Username))
		} else if update.CallbackQuery != nil {
			continue
		}
	}
	return errors.New("Polling closed")
}
