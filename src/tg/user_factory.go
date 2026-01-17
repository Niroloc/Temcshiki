package tg

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Operation string

const (
	ADD  Operation = "add"
	EDIT Operation = "edit"
)

type UserFactory struct {
	alias  string
	logger *logger.Logger
	args   UserFactoryArgs
}

type UserFactoryArgs struct {
	parsed    int
	operation Operation
	tgId      int
	username  string
	rights    data.UserRights
}

func CreateUserFactoryArgs() UserFactoryArgs {
	return UserFactoryArgs{
		parsed:    0,
		operation: "",
		tgId:      -1,
		username:  "",
	}
}

func CreateUserFactory() *UserFactory {
	return &UserFactory{
		alias:  "user",
		logger: logger.GetLogger(reflect.TypeFor[UserFactory]()),
		args:   CreateUserFactoryArgs(),
	}
}

func (this *UserFactory) GetAlias() string {
	return this.alias
}

func (this *UserFactory) ParseArguments(query *telego.CallbackQuery) error {
	args := strings.Split(query.Data, "_")[1:]
	this.args = CreateUserFactoryArgs()
	for i := 0; i < len(args); i++ {
		err := this.args.parseNext(args)
		if err != nil {
			this.logger.Error(err.Error())
			this.args = CreateUserFactoryArgs()
			break
		}
	}
	return nil
}

func (this *UserFactoryArgs) parseNext(args []string) error {
	if len(args) <= this.parsed {
		return errors.New("No next arg")
	}
	if this.parsed == 0 {
		this.operation = Operation(args[0])
		if this.operation != ADD && this.operation != EDIT {
			return errors.New("Unsupported operation")
		}
	} else if this.parsed == 1 {
		tgId, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Error while parsing tgID" + err.Error())
		}
		this.tgId = tgId
	} else if this.parsed == 2 {
		this.username = args[2]
	} else if this.parsed == 3 {
		this.rights = data.UserRights(args[3])
		if !data.RightsIsCorrect(this.rights) {
			return errors.New("Unsupported rights got")
		}
	} else {
		return errors.New("Unexpected arguments length value")
	}
	this.parsed++
	return nil
}

func (this *UserFactory) Apply(query *telego.CallbackQuery, user *data.User, bot *Bot) error {
	switch this.args.parsed {
	case 0:
		this.logger.Info("Stage 0 for user factory")
		markup := tu.InlineKeyboard(
			[]telego.InlineKeyboardButton{
				tu.InlineKeyboardButton("Добавить пользователя").WithCallbackData("user_add"),
				tu.InlineKeyboardButton("Редактировать роли").WithCallbackData("user_edit"),
			},
		)
		return bot.SendMessageWithMarkup(user, "Хм, что-то странное, но давай ещё раз попробуем", markup)
	case 1:
		switch this.args.operation {
		case ADD:
			err := bot.SendMessageWithoutMarkup(user, "Напишите сообщением tgID нового пользователя")
			if err != nil {
				return err
			}
			user.History.InputMode = true
			cbd := strings.Clone(query.Data)
			user.History.LastCallbackData = &cbd
			return nil
		case EDIT:
			return nil
		default:
			return nil
		}
	case 2:
		return nil
	case 3:
		return nil
	case 4:
		return nil
	default:
		return nil
	}
}
