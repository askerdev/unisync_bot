package tg

import (
	"net/http"
)

type Bot struct {
	Token  string
	Url    string
	Client *http.Client
}

func (b *Bot) url(method string) string {
	return b.Url + b.Token + "/" + method
}

type SendMessageParams struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (b *Bot) SendMessage(in *SendMessageParams) (*Message, error) {
	return Request[*Message](b.url("sendMessage"), in, b.Client)
}
