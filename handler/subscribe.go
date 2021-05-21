package handler

import (
	"errors"
	"regexp"

	"dotabot-ui/state"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

const (
	ASK_FOR_ACCOUNT_ID = "ASK_FOR_ACCOUNT_ID"
)

type Subscribe struct {
	subscribeRe *regexp.Regexp
	accountIdRe *regexp.Regexp

	repository repository.SubscriptionRepository
	telegramApi telegram.TelegramApi
	stateRepository state.StateRepository
}

func CreateSubscribe(repository repository.SubscriptionRepository,
						telegramApi telegram.TelegramApi,
						stateRepository state.StateRepository) *Subscribe {
	return &Subscribe{
		subscribeRe: regexp.MustCompile("^/?subscribe(@.*)?"),
		accountIdRe: regexp.MustCompile("\\d+"),
		repository: repository,
		telegramApi: telegramApi,
		stateRepository: stateRepository,
	}
}

func (this *Subscribe) CanHandle(update local_telegram.Update, state string) bool {
	switch(state) {
	case INIT_STATE:
		return this.subscribeRe.MatchString(update.Message.Text)
	case ASK_FOR_ACCOUNT_ID:
		return true
	default:
		return false
	}
}

func (this *Subscribe) Handle(update local_telegram.Update, state string) error {
	chat_id := update.Message.Chat.Id
	message := update.Message.Text

	switch(state) {
	case INIT_STATE:
		err := this.telegramApi.SendMessage(chat_id, "Enter account id")
		if err != nil {
			return err
		}
		err = this.stateRepository.SaveState(chat_id, ASK_FOR_ACCOUNT_ID)
		return err
	case ASK_FOR_ACCOUNT_ID:
		if !this.accountIdRe.MatchString(message) {
			err := this.telegramApi.SendMessage(chat_id, "Unrecognized account id")
			if err != nil {
				return err
			}
			return this.telegramApi.SendMessage(chat_id, "Enter account id")
		}

		matches := this.accountIdRe.FindStringSubmatch(message)
		if len(matches) != 1 {
			return errors.New("accountId was not found")
		}
		dotaAccountId := matches[0]
		
		subscription := repository.TelegramMatchSubscription {
			ChatId: chat_id,
			DotaAccountId: dotaAccountId,
		}
		err := this.repository.SaveLastKnownMatchId(subscription, 0)
		if err != nil {
			return err
		}
		return this.stateRepository.SaveState(chat_id, INIT_STATE)
	default:
		return errors.New("Unsupported state: " + state)
	}
}
