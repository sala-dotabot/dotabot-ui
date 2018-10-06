package repository

type MapRepository struct {
	holder map[TelegramMatchSubscription]uint64
}

func CreateMapRepository() MapRepository {
	return MapRepository {
		holder: make(map[TelegramMatchSubscription]uint64),
	}
}

func (this MapRepository) FindAll() ([]TelegramMatchSubscription, error) {
	var keys []TelegramMatchSubscription
	for k := range this.holder {
		keys = append(keys, k)
	}
	return keys, nil
}

func (this MapRepository) AddSubscription(subscription TelegramMatchSubscription) error {
	this.holder[subscription] = 0
	return nil
}

func (this MapRepository) RemoveSubscription(subscription TelegramMatchSubscription) error {
	delete(this.holder, subscription)
	return nil
}