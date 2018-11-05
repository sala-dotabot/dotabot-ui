package handler 

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	local_telegram "dotabot-ui/telegram"

	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

type Handler struct {
	subscribeRe *regexp.Regexp
	unsubscribeRe *regexp.Regexp
	repository repository.SubscriptionRepository
	telegramApi telegram.TelegramApi
}

func CreateHandler(repository repository.SubscriptionRepository, telegramApi telegram.TelegramApi) *Handler {
	return &Handler{
		repository: repository,
		telegramApi: telegramApi,
		subscribeRe: regexp.MustCompile("^/?subscribe (\\d+)"),
		unsubscribeRe: regexp.MustCompile("^/?unsubscribe (\\d+)"),
	}
}

func (this Handler) Handle(w http.ResponseWriter, r *http.Request) {
    if (r.URL == nil || r.URL.Path != "/webhook") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var u local_telegram.Update
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%#v", u.Message)

	text := u.Message.Text
	chatId := u.Message.Chat.Id

	log.Printf("Text: %s", text)

	if (text == "subscriptions" || text == "/subscriptions") {
		err = this.handleSubscriptions(chatId)
	} else if (this.subscribeRe.MatchString(text)) {
		err = this.handleSubscribe(chatId, text)
	} else if (this.unsubscribeRe.MatchString(text)) {
		err = this.handleUnscubscribe(chatId, text)
	} else {
		log.Printf("Handler was not found")
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}

func (this Handler) handleSubscriptions(chat_id int64) (err error) {
	log.Printf("Handle subscriptions %d", chat_id)
	repositories, err := this.repository.FindAll()

	if len(repositories) == 0 {
		this.telegramApi.SendMessage(chat_id, "No subscriptions yet!")
		return
	}

	str := ""
	for _, element := range repositories {
		if str != "" {
			str += "\n"
		}
		str += element.DotaAccountId
	}
	
	this.telegramApi.SendMessage(chat_id, str)
	return
}

func (this Handler) handleSubscribe(chat_id int64, message string) error {
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

func (this Handler) handleUnscubscribe(chat_id int64, message string) error {
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