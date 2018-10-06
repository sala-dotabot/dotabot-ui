package repository

type TelegramMatchSubscription struct {
	ChatId int64
	DotaAccountId string
}

type Repository interface {
	FindAll() (result []TelegramMatchSubscription, err error)
	AddSubscription(subscription TelegramMatchSubscription) error
	RemoveSubscription(subscription TelegramMatchSubscription) error
}