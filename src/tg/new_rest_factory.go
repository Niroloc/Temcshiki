package tg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
)

type NewRestFactory struct {
	alias  data.Alias
	args   *NewRestFactoryArgs
	logger *logger.Logger
}

type NewRestFactoryArgs struct {
	parsed       int
	name         string
	mapUrlId     int
	closestMetro string
}

func CreateNewRestFactoryArgs() *NewRestFactoryArgs {
	return &NewRestFactoryArgs{
		parsed:       0,
		name:         "",
		mapUrlId:     -1,
		closestMetro: "",
	}
}

func (this *NewRestFactoryArgs) parseNext(arg string) error {
	switch this.parsed {
	case 0:
		this.name = arg
	case 1:
		mapUrlId, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		this.mapUrlId = mapUrlId
	case 2:
		this.closestMetro = arg
	default:
		return errors.New("More than expected arguments for callback")
	}
	this.parsed++
	return nil
}

func CreateNewRestFactory() *NewRestFactory {
	return &NewRestFactory{
		alias:  data.NEW_REST,
		args:   nil,
		logger: logger.GetLogger(reflect.TypeFor[NewRestFactory]()),
	}
}

func (this *NewRestFactory) GetAlias() data.Alias {
	return this.alias
}

func (this *NewRestFactory) ParseArguments(query *telego.CallbackQuery) error {
	this.args = CreateNewRestFactoryArgs()
	for _, arg := range strings.Split(query.Data, "_")[1:] {
		if this.args.parseNext(arg) != nil {
			this.logger.Error(
				fmt.Sprintf(
					"Error while parsing argument %s as %d arg",
					arg,
					this.args.parsed,
				),
			)
			return nil
		}
	}
	return nil
}

func (this *NewRestFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	switch this.args.parsed {
	case 0:
		if err := bot.SendMessageWithoutMarkup(user, "Отправьте название ресторана"); err != nil {
			return err
		}
		continueWithInputMode(user, query.Data)
		return nil
	case 1:
		if err := bot.SendMessageWithoutMarkup(user, "Отправьте ссылку на ресторан/карты"); err != nil {
			return err
		}
		continueWithInputMode(user, query.Data)
		return nil
	case 2:
		_, err := exportedData.GetCommonData().GetUrl(this.args.mapUrlId)
		if err != nil {
			return err
		}
		if err = bot.SendMessageWithoutMarkup(user, "Отправьте ближайшее метро для генерации места встречи"); err != nil {
			return err
		}
		continueWithInputMode(user, query.Data)
		return nil
	case 3:
		mapUrl, err := exportedData.GetCommonData().GetUrl(this.args.mapUrlId)
		if err != nil {
			bot.SendMessageWithoutMarkup(user, "Ресторан не добавлен, там какая-то ошибка")
			return err
		}
		if err := exportedData.AddRest(this.args.name, mapUrl.Link, user.Id, this.args.closestMetro); err != nil {
			bot.SendMessageWithoutMarkup(user, "Ресторан не добавлен, там какая-то ошибка")
			return err
		}
		if err = bot.SendMessageWithoutMarkup(
			user,
			fmt.Sprintf(
				"Добавлен вариант ресторана с названием %s, ссылка: %s, ближайшее метро: %s",
				this.args.name,
				mapUrl.Link,
				this.args.closestMetro,
			),
		); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("Unexpected args number value")
	}
}
