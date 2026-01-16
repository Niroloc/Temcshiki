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

func GetDefaultKeyboard(user *data.User) *telego.InlineKeyboardMarkup {
	switch user.Rights {
	case data.ADMIN:
		return tu.InlineKeyboard()
	case data.RESERVATOR:
		return tu.InlineKeyboard()
	case data.VISITOR:
		return tu.InlineKeyboard()
	case data.SPECTATOR:
		return tu.InlineKeyboard()
	default:
		return tu.InlineKeyboard()
	}
}

type Bot struct {
	bot             *telego.Bot
	callbackManager *CallbackFactoryManager
	messageFactory  *MessageFactory
	logger          *logger.Logger
}

func CreateBot(botToken string) *Bot {
	logger := logger.GetLogger(reflect.TypeFor[Bot]())
	bot, err := telego.NewBot(botToken)
	if err != nil {
		logger.Error("Error while creating bot")
		panic(err)
	}
	return &Bot{
		bot:             bot,
		callbackManager: CreatCallbackFactoryManager([]CallbackFactory{}),
		logger:          logger,
	}
}

func (this *Bot) SendMessage(user *data.User, msg string) error {
	ctx := context.Background()
	this.logger.Info(fmt.Sprintf("Sending default message to %d", user.TgId))
	chatId := tu.ID(int64(user.TgId))
	_, err := this.bot.SendMessage(ctx, tu.Message(chatId, msg).WithReplyMarkup(GetDefaultKeyboard(user)))
	return err
}

func (this *Bot) SendMessageWithVoting(user *data.User, prefix string, eventId int, objs []data.VotingObject) error {
	this.logger.Info(fmt.Sprintf("Sending voting to %d", user.TgId))
	ctx := context.Background()
	chatId := tu.ID(int64(user.TgId))
	restsButtons := []telego.InlineKeyboardButton{}
	suffix := ""
	for i, o := range objs {
		suffix += o.GetDescription(i)
		restsButtons = append(
			restsButtons,
			tu.InlineKeyboardButton(o.GetButtonTitle()).WithCallbackData(o.GetCallbackData(eventId)),
		)
	}
	buttonRows := [][]telego.InlineKeyboardButton{}
	for i := 0; i < len(restsButtons); i += 3 {
		buttonRows = append(buttonRows, restsButtons[i:i+3])
	}
	inlineKeyboard := tu.InlineKeyboard(buttonRows...)
	_, err := this.bot.SendMessage(
		ctx,
		tu.Message(chatId, strings.Join([]string{prefix, suffix}, "\n")).WithReplyMarkup(inlineKeyboard),
	)
	return err
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
			user, err := exportedData.GetUser(tgId)
			if err != nil {
				this.logger.Warn(fmt.Sprintf("An unknown user is trying to use the bot! Id: %d", tgId))
				continue
			}
			this.SendMessage(user, fmt.Sprint(exportedData.GetUsers()[tgId].Username))
		} else if update.CallbackQuery != nil {
			continue
		}
	}
	return errors.New("Polling closed")
}
