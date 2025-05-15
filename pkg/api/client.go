package api

import (
	"github.com/go-resty/resty/v2"
)

func NewRestyClient(cfg Config) *resty.Client {
	return resty.New().
		SetBaseURL(cfg.HostURL).
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryMaxWaitTime(cfg.RetryMaxWaitTime).
		SetDebug(cfg.Debug)
}
