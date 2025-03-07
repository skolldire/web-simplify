package rest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/sony/gobreaker/v2"
)

var _ Service = (*client)(nil)

func NewClient(cfg Config, l log.Service) Service {
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

func checkBreakerState(counts gobreaker.Counts, c Config, l log.Service) bool {
	var failureRate float64
	if counts.Requests > 0 {
		failureRate = float64(counts.TotalFailures) / float64(counts.Requests)
	}
	l.Info(context.Background(), "Circuit Breaker Metrics",
		map[string]interface{}{
			"Total Requests":     counts.Requests,
			"Total Successes":    counts.TotalSuccesses,
			"Total Failures":     counts.TotalFailures,
			"Failure Rate":       failureRate,
			"ConsecutiveFails":   counts.ConsecutiveFailures,
			"ConsecutiveSuccess": counts.ConsecutiveSuccesses,
		})
	if counts.ConsecutiveFailures > c.CBMaxRequests || (counts.Requests >= c.CBRequestThreshold && failureRate > c.CBFailureRateLimit) {
		l.Info(context.Background(), "Circuit Breaker se abrirá debido a una alta tasa de fallos.", nil)
		return true
	}
	return false
}

func createCB(c Config, l log.Service) *gobreaker.CircuitBreaker[any] {
	if !c.WithCB {
		return nil
	}
	cbConfig := gobreaker.Settings{
		Name:        c.CBName,
		MaxRequests: c.CBMaxRequests,
		Interval:    c.CBInterval,
		Timeout:     c.CBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return checkBreakerState(counts, c, l)
		},

		OnStateChange: func(name string, from, to gobreaker.State) {
			l.Warn(context.Background(), map[string]interface{}{
				"name": name,
				"from": from,
				"to":   to,
			})
		},
	}

	cb := gobreaker.NewCircuitBreaker[any](cbConfig)
	return cb
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
			resp, err := reqFunc(ctx)

			if err != nil {
				c.logger.Warn(ctx, map[string]interface{}{
					"event": "request_failed",
					"error": err,
				})
				return nil, err
			}

			if resp != nil && resp.StatusCode() >= 500 {
				c.logger.Warn(ctx, map[string]interface{}{
					"event":  "server_error",
					"status": resp.StatusCode(),
				})
				return nil, fmt.Errorf("server error: HTTP %d", resp.StatusCode())
			}

			return resp, nil
		})

		if err != nil {
			if errors.Is(err, gobreaker.ErrOpenState) {
				c.logger.Error(ctx, err, "circuit breaker is open", map[string]interface{}{
					"event": "circuit_breaker_open",
					"error": err,
				})
				return nil, fmt.Errorf("circuit breaker open: %w", err)
			}
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
