package rest

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/sony/gobreaker"
)

type client struct {
	baseURL   string
	requester *resty.Request
	breaker   *gobreaker.CircuitBreaker
	tracer    log.Service
}

var _ Client = (*client)(nil)

func NewService(d Dependencies) *client {
	return &client{
		baseURL:   d.BaseURL,
		requester: d.Requester.HttpClient.R(),
		breaker:   d.Requester.Breaker,
		tracer:    d.Log,
	}
}

func (c *client) WithPathParams(pathParams map[string]string) *resty.Request {
	if pathParams != nil {
		c.requester.SetPathParams(pathParams)
	}
	return c.requester
}

func (c *client) WithHeaders(headers map[string]string) *resty.Request {
	if headers != nil {
		for key, value := range headers {
			c.requester.SetHeader(key, value)
		}
	}
	return c.requester
}

func (c *client) WithQueryParams(queryParams map[string]string) *resty.Request {
	if queryParams != nil {
		c.requester.SetQueryParams(queryParams)
	}
	return c.requester
}

func (c *client) WithBody(body interface{}) *resty.Request {
	if body != nil {
		c.requester.SetBody(body)
	}
	return c.requester
}

func (c *client) executeRequest(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	if c.breaker != nil {
		result, err := c.breaker.Execute(func() (interface{}, error) {
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
		resp, err := c.requester.SetContext(ctx).Get(c.baseURL + endpoint)
		if err != nil {
			c.tracer.Error(ctx, err, nil)
		}
		return resp, err
	})
}

func (c *client) Post(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		resp, err := c.requester.Post(c.baseURL + endpoint)
		if err != nil {
			c.tracer.Error(ctx, err, nil)
		}
		return resp, err
	})
}

func (c *client) Put(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		resp, err := c.requester.Put(c.baseURL + endpoint)
		if err != nil {
			c.tracer.Error(ctx, err, nil)
		}
		return resp, err
	})
}

func (c *client) Patch(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		resp, err := c.requester.Patch(c.baseURL + endpoint)
		if err != nil {
			c.tracer.Error(ctx, err, nil)
		}
		return resp, err
	})
}

func (c *client) Delete(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		resp, err := c.requester.Delete(c.baseURL + endpoint)
		if err != nil {
			c.tracer.Error(ctx, err, nil)
		}
		return resp, err
	})
}
