package requester

import (
	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"math"
	"math/rand"
	"time"
)

type Requester struct {
	HttpClient *resty.Client
	Breaker    *gobreaker.CircuitBreaker
}

func NewService(c Config) *Requester {
	return &Requester{
		HttpClient: createHttpClient(c),
		Breaker:    createCB(c),
	}
}

func createCB(c Config) *gobreaker.CircuitBreaker {
	if c.WithCB {
		return gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "HTTP client Circuit Breaker",
			MaxRequests: c.CBMaxRequests,
			Interval:    c.CBInterval * time.Second,
			Timeout:     c.CBTimeout * time.Second,
			ReadyToTrip: failureRateThreshold(c.CBRequestThreshold, c.CBFailureRateLimit),
		})
	}
	return nil
}

func createHttpClient(c Config) *resty.Client {
	if c.WithRetry {
		return resty.New().
			SetRetryCount(c.RetryCount).
			SetRetryAfter(RetryAfterFunc(c.RetryWaitTime*time.Millisecond, c.RetryMaxWaitTime*time.Second)).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})
	}
	return resty.New()
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

func ExponentialBackoffWithJitter(initialWaitTime, maxWaitTime time.Duration, attempt int) time.Duration {
	waitTime := initialWaitTime * time.Duration(math.Pow(2, float64(attempt-1)))
	jitter := time.Duration(rand.Float64() * float64(waitTime) * 0.5)
	waitTime += jitter

	if waitTime > maxWaitTime {
		waitTime = maxWaitTime
	}

	return waitTime
}

func RetryAfterFunc(initialWaitTime, maxWaitTime time.Duration) func(*resty.Client, *resty.Response) (time.Duration, error) {
	return func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		attempt := resp.Request.Attempt
		return ExponentialBackoffWithJitter(initialWaitTime, maxWaitTime, attempt), nil
	}
}
