package telegram

type Update struct {
	UpdateId int `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	MessageId int `json:"message_id"`
	Date int `json:"date"`
	Chat *Chat `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
}