package requester

import (
	"time"
)

type Config struct {
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
