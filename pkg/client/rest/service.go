package rest

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/sony/gobreaker"
	"math"
	"math/rand"
	"time"
)

var _ Service = (*client)(nil)

func NewService(cfg Config, l log.Service) Service {
	r := &requester{
		httpClient: createHttpClient(cfg, l),
		breaker:    createCB(cfg, l),
	}

	return &client{
		baseURL:   cfg.BaseURL,
		requester: r,
		logger:    l,
	}
}

func createCB(c Config, l log.Service) *gobreaker.CircuitBreaker {
	if !c.WithCB {
		return nil
	}

	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        c.CBName,
		MaxRequests: c.CBMaxRequests,
		Interval:    c.CBInterval,
		Timeout:     c.CBTimeout,
		ReadyToTrip: failureRateThreshold(c.CBRequestThreshold, c.CBFailureRateLimit),
		OnStateChange: func(name string, from, to gobreaker.State) {
			l.Warn(context.Background(), map[string]interface{}{
				"name": name,
				"from": from,
				"to":   to,
			})
		},
	})
}

func exponentialBackoffWithJitter(initialWaitTime, maxWaitTime time.Duration, attempt int, l log.Service) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	baseWaitTime := initialWaitTime * time.Duration(math.Pow(BackoffFactor, float64(attempt-1)))
	jitter := time.Duration(rand.Float64() * float64(baseWaitTime) * MaxJitterPercentage)
	waitTime := baseWaitTime + jitter

	if waitTime > maxWaitTime {
		waitTime = maxWaitTime
	}

	l.Debug(context.Background(), map[string]interface{}{
		"attempt":   attempt,
		"baseTime":  baseWaitTime,
		"jitter":    jitter,
		"waitTime":  waitTime,
		"maxWait":   maxWaitTime,
		"factor":    BackoffFactor,
		"jitterPct": MaxJitterPercentage,
	})

	return waitTime
}

func failureRateThreshold(requestThreshold int, failureRateLimit float64) func(gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		if counts.Requests < uint32(requestThreshold) {
			return false
		}
		failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
		return failureRate > failureRateLimit
	}
}

func retryAfterFunc(initialWaitTime, maxWaitTime time.Duration, l log.Service) func(*resty.Client, *resty.Response) (time.Duration, error) {
	return func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		attempt := resp.Request.Attempt
		return exponentialBackoffWithJitter(initialWaitTime, maxWaitTime, attempt, l), nil
	}
}

func createHttpClient(c Config, l log.Service) *resty.Client {
	client := resty.New()

	if c.WithRetry {
		client.SetRetryCount(c.RetryCount).
			SetRetryAfter(retryAfterFunc(c.RetryWaitTime, c.RetryMaxWaitTime, l)).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})
	}

	return client
}

func (c *client) executeRequest(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	if c.requester.breaker != nil {
		result, err := c.requester.breaker.Execute(func() (interface{}, error) {
			return reqFunc(ctx)
		})
		if err != nil {
			return nil, err
		}
		return result.(*resty.Response), nil
	}
	return reqFunc(ctx)
}

func (c *client) Get(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetContext(ctx).Get(c.baseURL + endpoint)
	})
}

func (c *client) Post(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Post(c.baseURL + endpoint)
	})
}

func (c *client) Put(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Put(c.baseURL + endpoint)
	})
}

func (c *client) Patch(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Patch(c.baseURL + endpoint)
	})
}

func (c *client) Delete(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetContext(ctx).Delete(c.baseURL + endpoint)
	})
}
