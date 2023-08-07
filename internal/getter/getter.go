package getter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	openapiclient "github.com/jadegopher/task/internal/api/jaicp/generated"
)

var (
	ErrInternal = fmt.Errorf("internal error occured")
)

type Filter struct {
	Filters    openapiclient.FiltersDto
	Pagination Pagination
}

type Pagination struct {
	Page int32
	Size int32
}

type BotsSessionGetter struct {
	client *openapiclient.APIClient
	token  string
}

func NewBotsSessionGetter(client *openapiclient.APIClient, token string) *BotsSessionGetter {
	return &BotsSessionGetter{
		client: client,
		token:  token,
	}
}

func (b *BotsSessionGetter) GetBotsSessions(ctx context.Context, filter Filter) ([]openapiclient.SessionInfo, error) {
	resp, r, err := b.client.SessionsApi.GetSessionDataByFilter(ctx, b.token).
		Page(filter.Pagination.Page).
		Size(filter.Pagination.Size).
		NeedToReturnSessionLabels(true).
		FiltersDto(filter.Filters).Execute()
	if err != nil {
		slog.Error("failed get bots' sessions", "err", err.Error())
		return nil, ErrInternal
	}
	if r.StatusCode != http.StatusOK {
		slog.Error("failed get bots' sessions: non ok status received")
		return nil, ErrInternal

	}

	return resp.Sessions, nil
}
