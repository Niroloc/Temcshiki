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

type ApproveFactory struct {
	alias  data.Alias
	args   *ApproveFactoryArgs
	logger *logger.Logger
}

type ApproveFactoryArgs struct {
	parsed  int
	eventId int
}

func CreateApproveFactoryArgs() *ApproveFactoryArgs {
	return &ApproveFactoryArgs{
		parsed:  0,
		eventId: -1,
	}
}

func (this *ApproveFactoryArgs) parseNext(arg string) error {
	switch this.parsed {
	case 0:
		if eid, err := strconv.Atoi(arg); err != nil {
			return err
		} else {
			this.eventId = eid
		}
	default:
		return errors.New("More than expected arguments for callback")
	}
	this.parsed++
	return nil
}

func CreateApproveFactory() *ApproveFactory {
	return &ApproveFactory{
		alias:  data.APPROVE,
		args:   nil,
		logger: logger.GetLogger(reflect.TypeFor[NewRestFactory]()),
	}
}

func (this *ApproveFactory) GetAlias() data.Alias {
	return this.alias
}

func (this *ApproveFactory) ParseArguments(query *telego.CallbackQuery) error {
	this.args = CreateApproveFactoryArgs()
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

func (this *ApproveFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	switch this.args.parsed {
	case 2:
		if err := exportedData.AddApproveVote(user, this.args.eventId); err != nil {
			return bot.SendMessage(user, "Во время учёта голоса произошла ошибка!")
		}
		return bot.SendMessage(user, "Ваш голос учтён!")
	default:
		return bot.SendMessage(user, "Ошибка!")
	}
}
