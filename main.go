package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	openapiclient "github.com/jadegopher/task/internal/api/jaicp/generated"

	"github.com/jadegopher/task/internal/entities"
	"github.com/jadegopher/task/internal/getter"
)

//Техническое задание:
//
//Имеется метод получения сессий бота по фильтру.
//Документация по методу с описанием можно найти по ссылке:
//https://app.jaicp.com/reporter/static/specs/reporter-public-en.html#tag/Sessions
//
//Необходимо написать функцию получения сессий бота за указанный интервал  времени и вернуть полученные данные в переменной соответствующего типа или ошибку.
//
//
//Для проверки полученных данных можно использовать
//{token} -
//В данном боте имеется 2 сессий за период 15 марта 2023 с 10:00:00 до 12:59:59 по UTC.
//
//Для фильтрации сессий по дате и времени нужно использовать тип фильтра DATE_TIME_RANGE
//
//Готовую функцию можно передать любым удобным способом, например, в https://go.dev/play/.

//go:embed config.json
var rawConfig []byte

func main() {
	config, err := entities.ParseConfig(rawConfig)
	if err != nil {
		slog.Error("failed parse config", "err", err.Error())
		return
	}

	sessionGetter := getter.NewBotsSessionGetter(
		openapiclient.NewAPIClient(openapiclient.NewConfiguration()),
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

	sessionGetter.GetBotsSessions(
		ctx, getter.Filter{
			Filters: openapiclient.FiltersDto{
				Filters: []openapiclient.AnalyticsFilter{
					{
						Key:  openapiclient.SESSION_ID,
						Type: "string",
					},
				},
			},
			Pagination: getter.Pagination{
				Page: 0,
				Size: 2,
			},
		},
	)

}
