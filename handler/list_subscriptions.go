package handler

import (
	"log"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

type ListSubscriptions struct {
	repository repository.SubscriptionRepository
	telegramApi telegram.TelegramApi
}

func CreateListSubscriptions(repository repository.SubscriptionRepository, 
								telegramApi telegram.TelegramApi) *ListSubscriptions {
	return &ListSubscriptions {
		repository: repository,
		telegramApi: telegramApi,
	}
}

func (this *ListSubscriptions) CanHandle(update local_telegram.Update) bool {
	text := update.Message.Text
	return text == "subscriptions" || text == "/subscriptions"
}

func (this *ListSubscriptions) Handle(update local_telegram.Update) error {
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