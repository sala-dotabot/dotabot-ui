package handler

import (
	"errors"
	"log"
	"regexp"

	"github.com/saladinkzn/dotabot-ui/state"

	local_telegram "github.com/saladinkzn/dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

const UNSUBSCRIBE_ASK_FOR_ACCOUNT_ID = "UNSUBSCRIBE_ASK_FOR_ACCOUNT_ID"

type Unsubscribe struct {
	accountIdRe   *regexp.Regexp
	unsubscribeRe *regexp.Regexp

	repository      repository.SubscriptionRepository
	stateRepository state.StateRepository
	telegramApi     telegram.TelegramApi
}

func CreateUnsubscribe(repository repository.SubscriptionRepository,
	telegramApi telegram.TelegramApi,
	stateRepository state.StateRepository) *Unsubscribe {
	return &Unsubscribe{
		accountIdRe:     regexp.MustCompile("\\d+"),
		unsubscribeRe:   regexp.MustCompile("^/?unsubscribe(@.*)?"),
		repository:      repository,
		stateRepository: stateRepository,
		telegramApi:     telegramApi,
	}
}

func (this *Unsubscribe) CanHandle(update local_telegram.Update, state string) bool {
	switch state {
	case INIT_STATE:
		return this.unsubscribeRe.MatchString(update.Message.Text)
	case UNSUBSCRIBE_ASK_FOR_ACCOUNT_ID:
		return true
	default:
		return false
	}
}

func (this *Unsubscribe) Handle(update local_telegram.Update, state string) error {
	chat_id := update.Message.Chat.Id
	message := update.Message.Text

	switch state {
	case INIT_STATE:
		err := this.telegramApi.SendMessage(chat_id, "Enter accountId")
		if err != nil {
			return err
		}
		return this.stateRepository.SaveState(chat_id, UNSUBSCRIBE_ASK_FOR_ACCOUNT_ID)
	case UNSUBSCRIBE_ASK_FOR_ACCOUNT_ID:
		log.Printf("Handle unsubscribe %d %s", chat_id, message)
		matches := this.accountIdRe.FindStringSubmatch(message)
		if len(matches) != 1 {
			return errors.New("accountId was not found")
		}
		dotaAccountId := matches[0]

		subscription := repository.TelegramMatchSubscription{
			ChatId:        chat_id,
			DotaAccountId: dotaAccountId,
		}
		err := this.repository.RemoveLastKnownMatchId(subscription)
		if err != nil {
			return err
		}
		return this.stateRepository.SaveState(chat_id, INIT_STATE)
	default:
		return nil
	}
}
