package cli

import (
	"context"
	"os"
	"time"

	"github.com/askerdev/unisync_bot/internal/converter"
	"github.com/askerdev/unisync_bot/internal/mospolytech"
	"github.com/askerdev/unisync_bot/internal/tg"
)

var GROUP = os.Getenv("GROUP")
var CHAT_ID = os.Getenv("CHAT_ID")

type Handler func(context.Context) error

type telegramBot interface {
	SendMessage(*tg.SendMessageParams) (*tg.Message, error)
}

type mospolytechAPI interface {
	Schedule() (*mospolytech.SemesterSchedule, error)
}

type cli struct {
	args     []string
	bot      telegramBot
	mpAPI    mospolytechAPI
	handlers map[string]Handler
}

func New(args []string, api mospolytechAPI, bot telegramBot) *cli {
	return &cli{
		args:     args,
		mpAPI:    api,
		bot:      bot,
		handlers: map[string]Handler{},
	}
}

func (a *cli) notify() error {
	sch, err := a.mpAPI.Schedule()
	if err != nil {
		return err
	}
	tasks, err := converter.TasksFromSchedule(CHAT_ID, GROUP, sch)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, t := range tasks {
		if t.TimeAt > now {
			continue
		}
		_, err := a.bot.SendMessage(&tg.SendMessageParams{
			ChatID:    t.ChatID,
			Text:      t.Text,
			ParseMode: "HTML",
		})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (a *cli) Run() error {
	return a.notify()
}
