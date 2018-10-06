package main

import (
	"dotabot-ui/handler"
	"dotabot-ui/telegram"

	"os"

	"github.com/saladinkzn/dotabot-cron/repository"
)

type Context struct {
	Handler *handler.Handler
}

func InitContext () (context *Context, err error) {
	telegramApiBaseUrl := getTelegramApiBaseUrl()
	telegramApiToken := getTelegramApiToken()
	telegramProxyUrl := getTelegramProxyUrl()

	telegramApi, err := telegram.CreateTelegramApiClient(telegramApiBaseUrl, telegramApiToken, telegramProxyUrl)
	if err != nil {
		return
	}

	repository := repository.CreateMapRepository()

	handler := handler.CreateHandler(repository, telegramApi)

	context = &Context {
		Handler: handler,
	}

	return
}


func getTelegramApiBaseUrl() string {
	telegramApiBaseUrl := os.Getenv("TELEGRAM_API_BASE_URL")
	if telegramApiBaseUrl != "" {
		return telegramApiBaseUrl
	} else {
		return "https://api.telegram.org"
	}
}

func getTelegramApiToken() string {
	return os.Getenv("TELEGRAM_API_TOKEN")
}

func getTelegramProxyUrl() string {
	return os.Getenv("TELEGRAM_PROXY_URL")
}