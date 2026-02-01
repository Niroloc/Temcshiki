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

type RestVoteFactory struct {
	alias  data.Alias
	args   *RestVoteFactoryArgs
	logger *logger.Logger
}

type RestVoteFactoryArgs struct {
	parsed  int
	eventId int
	restId  int
}

func CreateRestVoteFactoryArgs() *RestVoteFactoryArgs {
	return &RestVoteFactoryArgs{
		parsed:  0,
		eventId: -1,
		restId:  -1,
	}
}

func (this *RestVoteFactoryArgs) parseNext(arg string) error {
	switch this.parsed {
	case 0:
		if eid, err := strconv.Atoi(arg); err != nil {
			return err
		} else {
			this.eventId = eid
		}
	case 1:
		if rid, err := strconv.Atoi(arg); err != nil {
			return err
		} else {
			this.restId = rid
		}
	default:
		return errors.New("More than expected arguments for callback")
	}
	this.parsed++
	return nil
}

func CreateRestVoteFactory() *RestVoteFactory {
	return &RestVoteFactory{
		alias:  data.REST,
		args:   nil,
		logger: logger.GetLogger(reflect.TypeFor[NewRestFactory]()),
	}
}

func (this *RestVoteFactory) GetAlias() data.Alias {
	return this.alias
}

func (this *RestVoteFactory) ParseArguments(query *telego.CallbackQuery) error {
	this.args = CreateRestVoteFactoryArgs()
	for _, arg := range strings.Split(query.Data, "_")[1:] {
		if err := this.args.parseNext(arg); err != nil {
			this.logger.Error(
				fmt.Sprintf(
					"Error while parsing argument %s as %d arg",
					arg,
					this.args.parsed,
				),
			)
			return err
		}
	}
	return nil
}

func (this *RestVoteFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	switch this.args.parsed {
	case 2:
		if err := exportedData.AddRestVote(user, this.args.eventId, this.args.restId); err != nil {
			return bot.SendMessage(user, "Кажется, голос от вас за этот ресторан уже получен")
		}
		return bot.SendMessage(user, "Ваш голос учтён!")
	default:
		return bot.SendMessage(user, "Ошибка!")
	}
}
