package handler

import (
	"errors"
	"log"
	"regexp"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
)

type Subscribe struct {
	subscribeRe *regexp.Regexp

	repository repository.SubscriptionRepository
}

func CreateSubscribe(repository repository.SubscriptionRepository) *Subscribe {
	return &Subscribe{
		subscribeRe: regexp.MustCompile("^/?subscribe (\\d+)"),
		repository: repository,
	}
}

func (this *Subscribe) CanHandle(update local_telegram.Update) bool {
	return this.subscribeRe.MatchString(update.Message.Text)
}

func (this *Subscribe) Handle(update local_telegram.Update) error {
	chat_id := update.Message.Chat.Id
	message := update.Message.Text

	log.Printf("Handle subscribe %d %s", chat_id, message)
	matches := this.subscribeRe.FindStringSubmatch(message)
	if (len(matches) != 2) {
		return errors.New("accountId was not found")
	}
	dotaAccountId := matches[1]
	
	subscription := repository.TelegramMatchSubscription {
		ChatId: chat_id,
		DotaAccountId: dotaAccountId,
	}
	return this.repository.SaveLastKnownMatchId(subscription, 0)
}