package tg

type ID = int

type User struct {
	ID       ID     `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"is_bot"`
}

type Chat struct {
	ID   ID     `json:"id"`
	Type string `json:"type"`
}

type Message struct {
	MessageID  ID     `json:"message_id"`
	From       *User  `json:"from"`
	SenderChat *Chat  `json:"sender_chat"`
	Chat       *Chat  `json:"chat"`
	Text       string `json:"text"`
}
