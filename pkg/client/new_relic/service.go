package new_relic

import (
	"context"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
)

type client struct {
	app *newrelic.Application
	log *log.Service
}

var _ Service = (*client)(nil)

func NewClient(c Config) (*client, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigLicense(c.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(c.DistributedTracerEnabled),
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic: %w", err)
	}
	return &client{app: app}, nil
}

func (n *client) StartTransaction(ctx context.Context, name string) (context.Context, *newrelic.Transaction) {
	txn := n.app.StartTransaction(name)
	ctx = newrelic.NewContext(ctx, txn)
	return ctx, txn
}

func (n *client) EndTransaction(ctx context.Context) {
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		txn.End()
	}
}

func (n *client) RecordCustomEvent(ctx context.Context, eventType string, params map[string]interface{}) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.Application().RecordCustomEvent(eventType, params)
	}
}
