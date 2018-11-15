package state

type StateRepository interface {
	SaveState(chat_id int64, state string) error
	LoadState(chat_id int64) (string, error)
}

type MapRepository struct {
	holder map[int64]string
}

func CreateMapRepository() MapRepository {
	return MapRepository{
		holder: make(map[int64]string),
	}
}

func (this MapRepository) SaveState(chat_id int64, state string) error {
	this.holder[chat_id] = state
	return nil
}

func (this MapRepository) LoadState(chat_id int64) (string, error) {
	var err error = nil 
	var state = this.holder[chat_id]
	return state, err
}