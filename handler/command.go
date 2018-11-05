package handler

import (
	local_telegram "dotabot-ui/telegram"
)

type Command interface {
	CanHandle(update local_telegram.Update) bool

	Handle(update local_telegram.Update) error
}