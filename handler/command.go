package handler

import (
	local_telegram "github.com/saladinkzn/dotabot-ui/telegram"
)

type Command interface {
	CanHandle(update local_telegram.Update, state string) bool

	Handle(update local_telegram.Update, state string) error
}
