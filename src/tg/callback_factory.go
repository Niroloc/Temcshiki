package tg

import (
	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/mymmrac/telego"
)

type CallbackFactory interface {
	GetAlias() data.Alias
	ParseArguments(*telego.CallbackQuery) error
	Apply(*telego.CallbackQuery, *data.User, *data.Data, *Bot) error
}
