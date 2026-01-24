package tg

import (
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/mymmrac/telego"
)

type RatingFactory struct {
	alias  data.Alias
	logger *logger.Logger
}

func CreateRatingFactory() *RatingFactory {
	return &RatingFactory{
		alias:  data.RATING,
		logger: logger.GetLogger(reflect.TypeFor[UserFactory]()),
	}
}

func (this *RatingFactory) GetAlias() data.Alias {
	return this.alias
}

func (this *RatingFactory) ParseArguments(query *telego.CallbackQuery) error {
	return nil
}

func (this *RatingFactory) Apply(query *telego.CallbackQuery, user *data.User, exportedData *data.Data, bot *Bot) error {
	exportedData.GetRating()
	return nil
}
