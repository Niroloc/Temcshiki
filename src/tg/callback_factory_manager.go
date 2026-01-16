package tg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
)

type CallbackFactoryManager struct {
	aliasToFactory map[string]CallbackFactory
	logger         *logger.Logger
}

func CreatCallbackFactoryManager(cfs []CallbackFactory) *CallbackFactoryManager {
	aliasToFactory := map[string]CallbackFactory{}
	for _, cf := range cfs {
		aliasToFactory[cf.GetAlias()] = cf
	}
	logger := logger.GetLogger(reflect.TypeFor[CallbackFactoryManager]())
	return &CallbackFactoryManager{aliasToFactory: aliasToFactory, logger: logger}
}

func (this *CallbackFactoryManager) GetAndApplyFactory(callbackQuery *telego.CallbackQuery) {
	fact := this.aliasToFactory[strings.Split(callbackQuery.Data, "_")[0]]
	err := fact.ParseArguments(callbackQuery)
	if err != nil {
		this.logger.Error(
			fmt.Sprintf(
				"Parsing error for data: %s, factory: %s",
				callbackQuery.Data,
				reflect.TypeOf(fact),
			),
		)
	}
	err = fact.Apply(callbackQuery)
	if err != nil {
		this.logger.Error(
			fmt.Sprintf(
				"Applying error for data: %s, factory: %s",
				callbackQuery.Data,
				reflect.TypeOf(fact),
			),
		)
	}
}
