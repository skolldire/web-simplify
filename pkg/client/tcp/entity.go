package tcp

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net"
	"sync"
	"time"
)

var defaultConfig = Config{
	RetryCount:   5,
	RetryWait:    1 * time.Second,
	MaxRetryWait: 5 * time.Second,
	ConnTimeout:  3 * time.Second,
	ReadTimeout:  3 * time.Second,
	WriteTimeout: 3 * time.Second,
}

type Config struct {
	Address      string
	RetryCount   int
	RetryWait    time.Duration
	MaxRetryWait time.Duration
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Client interface {
	Connect() error
	SendMessage(msg string) (string, error)
	Close()
	HandleGracefulShutdown()
}

type client struct {
	log    log.Service
	config Config
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}
