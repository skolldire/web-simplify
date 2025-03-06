package rest

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/sony/gobreaker"
	"time"
)

const (
	BackoffFactor       = 2.0
	MaxJitterPercentage = 0.5
)

type Config struct {
	BaseURL            string
	WithRetry          bool
	RetryCount         int
	RetryWaitTime      time.Duration
	RetryMaxWaitTime   time.Duration
	WithCB             bool
	CBName             string
	CBMaxRequests      uint32
	CBInterval         time.Duration
	CBTimeout          time.Duration
	CBRequestThreshold int
	CBFailureRateLimit float64
}

type Service interface {
	Get(ctx context.Context, endpoint string) (*resty.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Patch(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Delete(ctx context.Context, endpoint string) (*resty.Response, error)
}

type client struct {
	baseURL   string
	requester *requester
	logger    log.Service
}

type requester struct {
	httpClient *resty.Client
	breaker    *gobreaker.CircuitBreaker
}
