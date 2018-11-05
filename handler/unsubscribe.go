package handler

import (
	"errors"
	"log"
	"regexp"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
)

type Unsubscribe struct {
	unsubscribeRe *regexp.Regexp
	repository repository.SubscriptionRepository
}

func CreateUnsubscribe(repository repository.SubscriptionRepository) *Unsubscribe {
	return &Unsubscribe{
		unsubscribeRe: regexp.MustCompile("^unsubscribe (\\d+)"),
		repository: repository,
	}
}

func (this *Unsubscribe) CanHandle(update local_telegram.Update) bool {
	return this.unsubscribeRe.MatchString(update.Message.Text)
}

func (this *Unsubscribe) Handle(update local_telegram.Update) error {
	chat_id := update.Message.Chat.Id 
	message := update.Message.Text

	log.Printf("Handle unsubscribe %d %s", chat_id, message)
	matches := this.unsubscribeRe.FindStringSubmatch(message)
	if (len(matches) != 2) {
		return errors.New("accountId was not found")
	}
	dotaAccountId := matches[1]

	subscription := repository.TelegramMatchSubscription {
		ChatId: chat_id,
		DotaAccountId: dotaAccountId,
	}
	return this.repository.RemoveLastKnownMatchId(subscription)
}
