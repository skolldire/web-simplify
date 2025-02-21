package rest

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/skolldire/web-simplify/pkg/utilities/requester"
)

type Client interface {
	Get(ctx context.Context, endpoint string) (*resty.Response, error)
	Post(ctx context.Context, endpoint string) (*resty.Response, error)
	Put(ctx context.Context, endpoint string) (*resty.Response, error)
	Patch(ctx context.Context, endpoint string) (*resty.Response, error)
	Delete(ctx context.Context, endpoint string) (*resty.Response, error)
}

type Dependencies struct {
	BaseURL   string
	Requester requester.Requester
	Log       log.Service
}

type RequestOptions struct {
	pathParams  map[string]string
	headers     map[string]string
	queryParams map[string]string
	body        interface{}
}
