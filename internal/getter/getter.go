package getter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	openapiclient "github.com/jadegopher/task/internal/api"
)

var (
	ErrInternal = fmt.Errorf("internal error occured")
)

type Filter struct {
	From       time.Time
	To         time.Time
	Pagination Pagination
}

type Pagination struct {
	Page int32
	Size int32
}

type BotsSessionGetter struct {
	client *openapiclient.Client
	token  string
}

func NewBotsSessionGetter(client *openapiclient.Client, token string) *BotsSessionGetter {
	return &BotsSessionGetter{
		client: client,
		token:  token,
	}
}

func (b *BotsSessionGetter) GetBotsSessions(ctx context.Context, filter Filter) ([]openapiclient.SessionInfo, error) {
	resp, err := b.client.GetSessionDataByFilter(
		ctx,
		b.token,
		&openapiclient.GetSessionDataByFilterParams{
			Page:                      &filter.Pagination.Page,
			Size:                      &filter.Pagination.Size,
			NeedToReturnSessionLabels: &[]bool{true}[0],
		},
		openapiclient.GetSessionDataByFilterJSONRequestBody{
			Filters: func() *[]openapiclient.AnalyticsFilter {
				return &[]openapiclient.AnalyticsFilter{
					{
						Key:  openapiclient.MESSAGETIME,
						Type: "DATE_TIME_RANGE",
						From: &filter.From,
						To:   &filter.To,
					},
				}
			}(),
		},
	)
	if err != nil {
		slog.Error("failed get bots' sessions", "err", err.Error())
		return nil, ErrInternal
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("failed get bots' sessions: non ok status received")
		return nil, ErrInternal
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed get bots' sessions: failed read body", "err", err.Error())
		return nil, ErrInternal
	}

	sessionData := &openapiclient.SessionsData{}
	if err = json.Unmarshal(data, sessionData); err != nil {
		slog.Error("failed get bots' sessions: failed unmarshall body", "err", err.Error())
		return nil, ErrInternal
	}

	return sessionData.Sessions, nil
}
