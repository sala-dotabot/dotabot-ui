package main

import (
	"dotabot-ui/handler"

	"os"

	"github.com/go-redis/redis"
	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
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

	client := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: "",
		DB:       0,
	})

	repository := repository.CreateRedisRepository(client)

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

func getRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}