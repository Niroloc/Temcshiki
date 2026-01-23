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

func getDefaultKeyboard(user *data.User) *telego.InlineKeyboardMarkup {
	switch user.Rights {
	case data.ADMIN:
		return getAdminKeyboard()
	case data.RESERVATOR:
		return getRegularUserKeyboard()
	case data.VISITOR:
		return getRegularUserKeyboard()
	case data.SPECTATOR:
		return getSpectatorKeyboard()
	default:
		return nil
	}
}

func getAdminKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Добавить пользователя").WithCallbackData("user_add"),
			tu.InlineKeyboardButton("Редактировать пользователя").WithCallbackData("user_edit"),
		},
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Получить рейтинг").WithCallbackData("rating"),
			tu.InlineKeyboardButton("Предложить ресторан").WithCallbackData("newrest"),
		},
	)
}

func getRegularUserKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Получить рейтинг").WithCallbackData("rating"),
			tu.InlineKeyboardButton("Предложить ресторан").WithCallbackData("newrest"),
		},
	)
}

func getSpectatorKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Получить рейтинг").WithCallbackData("rating"),
		},
	)
}

type Bot struct {
	bot             *telego.Bot
	callbackManager *CallbackFactoryManager
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
		bot: bot,
		callbackManager: CreateCallbackFactoryManager(
			[]CallbackFactory{
				CreateUserFactory(),
				CreateRatingFactory(),
			},
		),
		logger: logger,
	}
}

func (this *Bot) SendMessageWithoutMarkup(user *data.User, msg string) error {
	ctx := context.Background()
	this.logger.Info(fmt.Sprintf("Sending message without markup to %s", user.Username))
	chatId := tu.ID(int64(user.TgId))
	_, err := this.bot.SendMessage(ctx, tu.Message(chatId, msg))
	return err
}

func (this *Bot) SendMessage(user *data.User, msg string) error {
	ctx := context.Background()
	this.logger.Info(fmt.Sprintf("Sending default message to %s", user.Username))
	chatId := tu.ID(int64(user.TgId))
	replyMarkup := getDefaultKeyboard(user)
	var err error
	if replyMarkup != nil {
		_, err = this.bot.SendMessage(ctx, tu.Message(chatId, msg).WithReplyMarkup(replyMarkup))
	} else {
		_, err = this.bot.SendMessage(ctx, tu.Message(chatId, msg))
	}
	return err
}

func (this *Bot) SendMessageWithVoting(user *data.User, prefix string, eventId int, objs []data.VotingObject) error {
	this.logger.Info(fmt.Sprintf("Sending voting to %s", user.Username))
	ctx := context.Background()
	chatId := tu.ID(int64(user.TgId))
	buts := []telego.InlineKeyboardButton{}
	suffix := ""
	for i, o := range objs {
		suffix += o.GetDescription(i)
		buts = append(
			buts,
			tu.InlineKeyboardButton(o.GetButtonTitle()).WithCallbackData(o.GetCallbackData(eventId)),
		)
	}
	buttonRows := [][]telego.InlineKeyboardButton{}
	for i := 0; i < len(buts); i += 3 {
		buttonRows = append(buttonRows, buts[i:min(len(buts), i+3)])
	}
	inlineKeyboard := tu.InlineKeyboard(buttonRows...)
	_, err := this.bot.SendMessage(
		ctx,
		tu.Message(chatId, strings.Join([]string{prefix, suffix}, "\n")).WithReplyMarkup(inlineKeyboard),
	)
	return err
}

func (this *Bot) SendMessageWithMarkup(user *data.User, msg string, markup *telego.InlineKeyboardMarkup) error {
	this.logger.Info(fmt.Sprintf("Sending markup message to %s", user.Username))
	ctx := context.Background()
	chatId := tu.ID(int64(user.TgId))
	_, err := this.bot.SendMessage(ctx, tu.Message(chatId, msg).WithReplyMarkup(markup))
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
			var user *data.User
			if user, err = this.welcome(update.Message, exportedData); err != nil {
				continue
			}
			if user.History.InputMode {
				markup := tu.InlineKeyboard(
					[]telego.InlineKeyboardButton{
						tu.InlineKeyboardButton("Отмена").WithCallbackData(*user.History.LastCallbackData),
						tu.InlineKeyboardButton("Продолжить").WithCallbackData(
							*user.History.LastCallbackData + "_" + update.Message.Text,
						),
					},
				)
				err := this.SendMessageWithMarkup(user, "Подтвердите корректность ввода", markup)
				if err != nil {
					this.logger.Error(err.Error())
				}
			} else {
				this.SendMessage(
					user,
					fmt.Sprintf(
						"Добро пожаловать в бота Темщиков!\n"+
							"Ваша роль: %s.\n"+
							"Для действий нажмите на одну из кнопок ниже.",
						user.Rights,
					),
				)
			}
		} else if update.CallbackQuery != nil {
			this.callbackManager.GetAndApplyFactory(update.CallbackQuery, exportedData, this)
		}
	}
	return errors.New("Polling closed")
}

func (this *Bot) welcome(msg *telego.Message, exportedData *data.Data) (*data.User, error) {
	tgId := msg.Chat.ID
	this.logger.Info(fmt.Sprintf("User with id %d has written to the bot", tgId))
	user, err := exportedData.GetUserByTg(int(tgId))
	if err != nil {
		this.logger.Warn(fmt.Sprintf("An unknown user has written to the bot. ID: %d", tgId))
		return nil, errors.New("Unknown user")
	}
	return user, nil
}
