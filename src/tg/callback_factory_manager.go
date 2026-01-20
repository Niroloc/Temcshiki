package tg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
)

type CallbackFactoryManager struct {
	aliasToFactory map[string]CallbackFactory
	logger         *logger.Logger
}

func CreateCallbackFactoryManager(cfs []CallbackFactory) *CallbackFactoryManager {
	aliasToFactory := map[string]CallbackFactory{}
	for _, cf := range cfs {
		aliasToFactory[cf.GetAlias()] = cf
	}
	logger := logger.GetLogger(reflect.TypeFor[CallbackFactoryManager]())
	return &CallbackFactoryManager{aliasToFactory: aliasToFactory, logger: logger}
}

func (this *CallbackFactoryManager) GetAndApplyFactory(callbackQuery *telego.CallbackQuery, exportedData *data.Data, bot *Bot) {
	user, err := exportedData.GetUserByTg(int(callbackQuery.Message.GetChat().ID))
	if err != nil {
		this.logger.Error(fmt.Sprintf("Unknown user is sending callback query. ID: %d", callbackQuery.Message.GetChat().ID))
		return
	}
	fact, exists := this.aliasToFactory[strings.Split(callbackQuery.Data, "_")[0]]
	if !exists {
		this.logger.Error(fmt.Sprintf("Unknown callback factory for data: %s", callbackQuery.Data))
		return
	}
	user.History.InputMode = false
	user.History.LastCallbackData = nil
	err = fact.ParseArguments(callbackQuery)
	if err != nil {
		this.logger.Error(
			fmt.Sprintf(
				"Parsing error for data: %s, factory: %s\nError: %s",
				callbackQuery.Data,
				reflect.TypeOf(fact),
				err.Error(),
			),
		)
		return
	}
	err = fact.Apply(callbackQuery, user, exportedData, bot)
	if err != nil {
		this.logger.Error(
			fmt.Sprintf(
				"Applying error for data: %s, factory: %s\nError: %s",
				callbackQuery.Data,
				reflect.TypeOf(fact),
				err.Error(),
			),
		)
	}
}
