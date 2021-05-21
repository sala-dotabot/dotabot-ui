package handler

import (
	"log"
	"regexp"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

type ListSubscriptions struct {
	subscriptionsRe *regexp.Regexp

	repository repository.SubscriptionRepository
	telegramApi telegram.TelegramApi
}

func CreateListSubscriptions(repository repository.SubscriptionRepository, 
								telegramApi telegram.TelegramApi) *ListSubscriptions {
	return &ListSubscriptions {
		subscriptionsRe: regexp.MustCompile("^/?subscriptions(@.*)?"),
		repository: repository,
		telegramApi: telegramApi,
	}
}

func (this *ListSubscriptions) CanHandle(update local_telegram.Update, state string) bool {
	if (state != INIT_STATE) {
		return false
	}

	text := update.Message.Text
	return this.subscriptionsRe.MatchString(text)
}

func (this *ListSubscriptions) Handle(update local_telegram.Update, state string) error {
	chat_id := update.Message.Chat.Id
	log.Printf("Handle subscriptions %d", chat_id)
	repositories, err := this.repository.FindAll()

	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		err := this.telegramApi.SendMessage(chat_id, "No subscriptions yet!")
		return err
	}

	str := ""
	for _, element := range repositories {
		if str != "" {
			str += "\n"
		}
		str += element.DotaAccountId
	}
	
	err = this.telegramApi.SendMessage(chat_id, str)
	return err
}