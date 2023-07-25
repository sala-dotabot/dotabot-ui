package state

import (
	"strconv"

	"github.com/go-redis/redis"
)

const fieldName = "state"

type RedisStateRepostory struct {
	client *redis.Client
}

func CreateRedisStateRepostory(client *redis.Client) StateRepository {
	return RedisStateRepostory{
		client: client,
	}
}

func (this RedisStateRepostory) SaveState(chat_id int64, state string) error {
	key := makeTelegramChatKey(chat_id)

	_, err := this.client.HSet(key, fieldName, state).Result()
	return err
}

func (this RedisStateRepostory) LoadState(chat_id int64) (result string, err error) {
	key := makeTelegramChatKey(chat_id)

	result, err = this.client.HGet(key, fieldName).Result()
	if err != nil {
		return
	}

	return
}

func makeTelegramChatKey(chat_id int64) string {
	return "telegramChatKey:" + strconv.FormatInt(chat_id, 10)
}
