package tg

import "github.com/mymmrac/telego"

type CallbackFactory interface {
	GetAlias() string
	ParseArguments(*telego.CallbackQuery) error
	Apply(*telego.CallbackQuery) error
}
