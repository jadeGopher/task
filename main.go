package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	openapiclient "github.com/jadegopher/task/internal/api"
	"github.com/jadegopher/task/internal/entities"
	"github.com/jadegopher/task/internal/getter"
)

//go:embed config.json
var rawConfig []byte

func main() {
	config, err := entities.ParseConfig(rawConfig)
	if err != nil {
		slog.Error("failed parse config", "err", err.Error())
		return
	}

	client, err := openapiclient.NewClient("https://app.jaicp.com")
	if err != nil {
		slog.Error("failed init client", "err", err.Error())
		return
	}

	sessionGetter := getter.NewBotsSessionGetter(
		client,
		config.Token,
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh)
		for {
			sig := <-signalCh
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				slog.Info("signal received. shutting down app", "signal", sig)
				cancelFunc()
			}
		}
	}()

	sessions, err := sessionGetter.GetBotsSessions(
		ctx, getter.Filter{
			From: time.Date(2023, 03, 15, 10, 00, 00, 00, time.UTC),
			To:   time.Date(2023, 03, 15, 12, 59, 59, 00, time.UTC),
			Pagination: getter.Pagination{
				Page: 0,
				Size: 2,
			},
		},
	)
	if err != nil {
		return
	}

	for _, session := range sessions {
		fmt.Println(session)
	}
}
