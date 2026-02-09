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

type DateVoteFactory struct {
	alias  data.Alias
	args   *DateVoteFactoryArgs
	logger *logger.Logger
}

type DateVoteFactoryArgs struct {
	parsed  int
	eventId int
	dateId  int
}

func CreateDateVoteFactoryArgs() *DateVoteFactoryArgs {
	return &DateVoteFactoryArgs{
		parsed:  0,
		eventId: -1,
		dateId:  -1,
	}
}

func (this *DateVoteFactoryArgs) parseNext(arg string) error {
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
			this.dateId = rid
		}
	default:
		return errors.New("More than expected arguments for callback")
	}
	this.parsed++
	return nil
}

func CreateDateVoteFactory() *DateVoteFactory {
	return &DateVoteFactory{
		alias:  data.DATE,
		args:   nil,
		logger: logger.GetLogger(reflect.TypeFor[NewRestFactory]()),
	}
}

func (this *DateVoteFactory) GetAlias() data.Alias {
	return this.alias
}

func (this *DateVoteFactory) ParseArguments(query *telego.CallbackQuery) error {
	this.args = CreateDateVoteFactoryArgs()
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

func (this *DateVoteFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	switch this.args.parsed {
	case 2:
		if err := exportedData.AddDateVote(user, this.args.eventId, this.args.dateId); err != nil {
			return bot.SendMessage(user, "Во время учёта голоса произошла ошибка!")
		}
		return bot.SendMessage(user, "Ваш голос учтён!")
	default:
		return bot.SendMessage(user, "Ошибка!")
	}
}
