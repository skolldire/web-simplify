package tcp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewClient(cfg Config, log log.Service) Client {
	ctx, cancel := context.WithCancel(context.Background())
	cfg = mergeWithDefaults(cfg, log)

	return &client{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
		log:    log,
	}
}

func mergeWithDefaults(cfg Config, log log.Service) Config {
	if cfg.Address == "" {
		log.FatalError(context.Background(), errors.New("missing required Address field in TCP client config"), nil)
	}

	if cfg.RetryCount == 0 {
		cfg.RetryCount = defaultConfig.RetryCount
	}
	if cfg.RetryWait == 0 {
		cfg.RetryWait = defaultConfig.RetryWait
	}
	if cfg.MaxRetryWait == 0 {
		cfg.MaxRetryWait = defaultConfig.MaxRetryWait
	}
	if cfg.ConnTimeout == 0 {
		cfg.ConnTimeout = defaultConfig.ConnTimeout
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = defaultConfig.ReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = defaultConfig.WriteTimeout
	}

	return cfg
}

func (c *client) Connect() error {
	var err error
	for i := 0; i < c.config.RetryCount; i++ {
		c.conn, err = net.DialTimeout("tcp", c.config.Address, c.config.ConnTimeout)
		if err == nil {
			c.log.Info(context.Background(), fmt.Sprintf("Conected..."), map[string]interface{}{"address": c.config.Address})
			return nil
		}
		waitTime := c.config.RetryWait * time.Duration(i+1)
		if waitTime > c.config.MaxRetryWait {
			waitTime = c.config.MaxRetryWait
		}
		c.log.Debug(context.Background(), map[string]interface{}{"retry": i + 1, "waitTime": waitTime})
		time.Sleep(waitTime)
	}
	return c.log.WrapError(err, fmt.Sprintf("failed to connect after %d.", c.config.RetryCount))
}

func (c *client) SendMessage(msg string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return "", c.log.WrapError(errors.New("connection is nil"), "error to send message")
	}

	_ = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))

	_, err := c.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return "", c.log.WrapError(err, "error writing message")
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))

	reader := bufio.NewReader(c.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", c.log.WrapError(err, "error reading response")
	}

	return response, nil
}

func (c *client) Close() {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
		fmt.Println("Conexión cerrada.")
	}
}

func (c *client) HandleGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		c.log.Info(context.Background(), "Shutdown signal received. Closing connection…", nil)
		c.Close()
		os.Exit(0)
	}()
}
