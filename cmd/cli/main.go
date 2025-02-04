package main

import (
	"context"
	_ "embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/askerdev/unisync_bot/internal/cli"
	"github.com/askerdev/unisync_bot/internal/mospolytech"
	"github.com/askerdev/unisync_bot/internal/tg"
	_ "github.com/mattn/go-sqlite3"
)

const TOKEN_KEY = "TELEGRAM_BOT_API_TOKEN"
const URL = "https://api.telegram.org/bot"
const SCHEDULE_URL = "https://rasp.dmami.ru/semester.json"

func main() {
	httpclient := &http.Client{}

	bot := &tg.Bot{
		Token:  os.Getenv(TOKEN_KEY),
		Url:    URL,
		Client: httpclient,
	}

	api := mospolytech.NewAPI(
		SCHEDULE_URL,
		httpclient,
	)

	app := cli.New(os.Args, api, bot)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	go func() {
		exit := make(
			chan os.Signal,
			1,
		)
		defer close(exit)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		select {
		case <-exit:
			slog.Warn("Exiting program! Recieved signal!")
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
