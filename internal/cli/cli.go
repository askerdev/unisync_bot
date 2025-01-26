package cli

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/askerdev/unisync_bot/internal/converter"
	"github.com/askerdev/unisync_bot/internal/domain"
	"github.com/askerdev/unisync_bot/internal/mospolytech"
	"github.com/askerdev/unisync_bot/internal/tg"
)

var GROUP = os.Getenv("GROUP")
var CHAT_ID = os.Getenv("CHAT_ID")

type Handler func(context.Context) error

type Storage interface {
	Insert(context.Context, []*domain.Task) error
	Select(context.Context) ([]*domain.Task, error)
	Delete(context.Context, int) error
}

type TelegramBot interface {
	SendMessage(*tg.SendMessageParams) (*tg.Message, error)
}

type MospolytechAPI interface {
	Schedule() (*mospolytech.SemesterSchedule, error)
}

type app struct {
	args     []string
	storage  Storage
	bot      TelegramBot
	mpAPI    MospolytechAPI
	handlers map[string]Handler
}

func NewApp(args []string, storage Storage, api MospolytechAPI, bot TelegramBot) *app {
	return &app{
		args:     args,
		storage:  storage,
		mpAPI:    api,
		bot:      bot,
		handlers: map[string]Handler{},
	}
}

func (a *app) test(ctx context.Context) error {
	tasks, err := a.storage.Select(ctx)
	if err != nil {
		return err
	}

	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	t := tasks[0]

	wg.Add(1)
	go func() {
		m.Lock()
		defer m.Unlock()
		defer wg.Done()
		a.storage.Delete(ctx, t.ID)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.bot.SendMessage(&tg.SendMessageParams{
			ChatID:    t.ChatID,
			Text:      t.Text,
			ParseMode: "HTML",
		})
	}()

	wg.Wait()

	return nil
}

func (a *app) update(ctx context.Context) error {
	sch, err := a.mpAPI.Schedule()
	if err != nil {
		return err
	}
	tasks, err := converter.TasksFromSchedule(CHAT_ID, GROUP, sch)
	if err != nil {
		return err
	}

	return a.storage.Insert(ctx, tasks)
}

func (a *app) notify(ctx context.Context) error {
	tasks, err := a.storage.Select(ctx)
	if err != nil {
		return err
	}

	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	now := time.Now().Unix()
	for _, t := range tasks {
		if t.TimeAt > now {
			continue
		}
		wg.Add(1)
		go func() {
			m.Lock()
			defer m.Unlock()
			defer wg.Done()
			a.storage.Delete(ctx, t.ID)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.bot.SendMessage(&tg.SendMessageParams{
				ChatID:    t.ChatID,
				Text:      t.Text,
				ParseMode: "HTML",
			})
		}()
	}

	wg.Wait()

	return nil
}

func (a *app) Run(ctx context.Context) error {
	a.handlers["update"] = a.update
	a.handlers["notify"] = a.notify
	a.handlers["test"] = a.test
	action, ok := a.handlers[a.args[1]]
	if !ok {
		return errors.New("cli: unknown command")
	}
	return action(ctx)
}
