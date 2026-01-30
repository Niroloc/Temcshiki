package tg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/utils"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Operation string

const (
	ADD  Operation = "add"
	EDIT Operation = "edit"
)

type UserFactory struct {
	alias  data.Alias
	logger *logger.Logger
	args   UserFactoryArgs
}

type UserFactoryArgs struct {
	parsed    int
	operation Operation
	id        int
	username  string
	rights    data.UserRights
}

func CreateUserFactoryArgs() UserFactoryArgs {
	return UserFactoryArgs{
		parsed:    0,
		operation: "",
		id:        -1,
		username:  "",
	}
}

func CreateUserFactory() *UserFactory {
	return &UserFactory{
		alias:  data.USER,
		logger: logger.GetLogger(reflect.TypeFor[UserFactory]()),
		args:   CreateUserFactoryArgs(),
	}
}

func (this *UserFactory) GetAlias() data.Alias {
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
			return err
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
		this.id = tgId
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

func (this *UserFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	switch this.args.parsed {
	case 0:
		this.logger.Info("Stage 0 for user factory")
		markup := tu.InlineKeyboard(
			[]telego.InlineKeyboardButton{
				tu.InlineKeyboardButton("Добавить пользователя").WithCallbackData(string(data.USER) + "_" + string(ADD)),
				tu.InlineKeyboardButton("Редактировать пользователя").WithCallbackData(string(data.USER) + "_" + string(EDIT)),
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
			continueWithInputMode(user, query.Data)
			return nil
		case EDIT:
			buts := utils.MapValuesToList(
				exportedData.GetUsers(),
				func(u *data.User) telego.InlineKeyboardButton {
					return tu.InlineKeyboardButton(u.Username).WithCallbackData(query.Data + fmt.Sprintf("_%d", u.Id))
				},
			)
			rows := [][]telego.InlineKeyboardButton{}
			for i := 0; i < len(buts); i += 3 {
				rows = append(rows, buts[i:min(i+3, len(buts))])
			}
			return bot.SendMessageWithMarkup(user, "Выберите пользователя, для изменения", tu.InlineKeyboard(rows...))
		default:
			return bot.SendMessage(user, "Что-то пошло не так, возвращаемся в главное меню")
		}
	case 2:
		switch this.args.operation {
		case ADD:
			err := bot.SendMessageWithoutMarkup(user, "Отлично, теперь напишите имя пользователя")
			if err != nil {
				return err
			}
			continueWithInputMode(user, query.Data)
			return nil
		case EDIT:
			editedUser, err := exportedData.GetUserById(this.args.id)
			if err != nil {
				return err
			}
			err = bot.SendMessageWithMarkup(
				user,
				"Напишите новое имя пользователя или нажмите на кнопку",
				tu.InlineKeyboard([]telego.InlineKeyboardButton{
					tu.InlineKeyboardButton("Оставить без изменений").WithCallbackData(query.Data + "_" + editedUser.Username),
				}),
			)
			if err != nil {
				return err
			}
			continueWithInputMode(user, query.Data)
			return nil
		default:
			return bot.SendMessage(user, "Что-то пошло не так, возвращаемся в главное меню")
		}
	case 3:
		rows := [][]telego.InlineKeyboardButton{
			{
				tu.InlineKeyboardButton("Админ").WithCallbackData(query.Data + "_" + string(data.ADMIN)),
				tu.InlineKeyboardButton("Резерватор").WithCallbackData(query.Data + "_" + string(data.RESERVATOR)),
			},
			{
				tu.InlineKeyboardButton("Посетитель").WithCallbackData(query.Data + "_" + string(data.VISITOR)),
				tu.InlineKeyboardButton("Наблюдатель").WithCallbackData(query.Data + "_" + string(data.SPECTATOR)),
			},
		}
		if this.args.operation == EDIT {
			editedUser, err := exportedData.GetUserById(this.args.id)
			if err != nil {
				return err
			}
			rows = append(rows, []telego.InlineKeyboardButton{tu.InlineKeyboardButton("Оставляем как есть").WithCallbackData(query.Data + "_" + editedUser.Username)})
		}
		return bot.SendMessageWithMarkup(
			user,
			"Теперь выбираем права нашего героя",
			tu.InlineKeyboard(rows...),
		)
	case 4:
		switch this.args.operation {
		case ADD:
			user := exportedData.CreateNewUser(this.args.id, this.args.username, this.args.rights)
			if user == nil {
				return bot.SendMessage(user, "Мы не смогли добавить пользователя, разбираемся...")
			}
		case EDIT:
			editedUser, err := exportedData.GetUserById(this.args.id)
			if err != nil {
				return bot.SendMessage(user, "Мы не смогли найти пользователя для редактирования, разбираемся...")
			}
			exportedData.UpdateUser(editedUser, this.args.username, this.args.rights)
		}
		return bot.SendMessage(user, "Всё готово, шеф, принимаем работу.")
	default:
		return nil
	}
}

func continueWithInputMode(user *data.User, cbd string) {
	user.History.InputMode = true
	user.History.LastCallbackData = &cbd
}
