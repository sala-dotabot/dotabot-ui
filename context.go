package main

import (
	"github.com/saladinkzn/dotabot-ui/handler"
	"github.com/saladinkzn/dotabot-ui/state"

	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

type Context struct {
	Handler *handler.Handler

	MetricsHandler http.Handler
}

func InitContext() (context *Context, err error) {
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
	stateRepository := state.CreateMapRepository()

	listSubscriptions := handler.CreateListSubscriptions(repository, telegramApi)
	subscribe := handler.CreateSubscribe(repository, telegramApi, stateRepository)
	unsubscribe := handler.CreateUnsubscribe(repository, telegramApi, stateRepository)

	commands := []handler.Command{
		listSubscriptions,
		subscribe,
		unsubscribe,
	}

	totalCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_update_count",
	})

	notFoundErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "not_found_error_count",
	})

	decodeErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "decode_error_count",
	})

	loadStateErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "load_state_error_count",
	})

	unhandledCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "unhandled_count",
	})

	errorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "error_count",
	})

	handler := handler.CreateHandler(commands, stateRepository, totalCount, notFoundErrorCounter, decodeErrorCounter, loadStateErrorCounter, unhandledCounter, errorCounter)
	metricsHandler := promhttp.Handler()

	context = &Context{
		Handler:        handler,
		MetricsHandler: metricsHandler,
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
