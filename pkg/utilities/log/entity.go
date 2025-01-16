package log

import "context"

type Service interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, err error, fields map[string]interface{})
	Debug(ctx context.Context, fields map[string]interface{})
	Warn(ctx context.Context, fields map[string]interface{})
	WrapError(err error, msg string) error
}
