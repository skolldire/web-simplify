package new_relic

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	StartTransaction(ctx context.Context, name string) (context.Context, *newrelic.Transaction)
	EndTransaction(ctx context.Context)
	RecordCustomEvent(ctx context.Context, eventType string, params map[string]interface{})
}

type Config struct {
	AppName                  string `json:"app_name"`
	LicenseKey               string `json:"license_key"`
	DistributedTracerEnabled bool   `json:"distributed_tracer_enabled"`
}
