package tg

import (
	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/mymmrac/telego"
)

type CallbackFactory interface {
	GetAlias() string
	ParseArguments(*telego.CallbackQuery) error
	Apply(*telego.CallbackQuery, *data.Data) error
}
