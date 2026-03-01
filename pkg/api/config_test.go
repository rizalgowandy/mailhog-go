package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfigValidateDefaults(t *testing.T) {
	cfg := Config{}

	err := cfg.Validate()
	require.NoError(t, err)
	require.Equal(t, time.Minute, cfg.Timeout)
	require.Equal(t, 2*time.Second, cfg.RetryMaxWaitTime)
	require.Equal(t, "http://localhost:8025", cfg.HostURL)
	require.Equal(t, 0, cfg.RetryCount)
}

func TestConfigValidateNegativeRetryCount(t *testing.T) {
	cfg := Config{RetryCount: -1}

	err := cfg.Validate()
	require.NoError(t, err)
	require.Equal(t, 3, cfg.RetryCount)
}

func TestConfigValidateKeepsProvidedValues(t *testing.T) {
	cfg := Config{
		Timeout:          5 * time.Second,
		RetryCount:       2,
		RetryMaxWaitTime: 7 * time.Second,
		Debug:            true,
		HostURL:          "http://example.test",
	}

	err := cfg.Validate()
	require.NoError(t, err)
	require.Equal(t, 5*time.Second, cfg.Timeout)
	require.Equal(t, 2, cfg.RetryCount)
	require.Equal(t, 7*time.Second, cfg.RetryMaxWaitTime)
	require.True(t, cfg.Debug)
	require.Equal(t, "http://example.test", cfg.HostURL)
}

func TestNewRestyClient(t *testing.T) {
	cfg := Config{
		Timeout:          3 * time.Second,
		RetryCount:       4,
		RetryMaxWaitTime: 2 * time.Second,
		Debug:            true,
		HostURL:          "http://localhost:9999",
	}

	cli := NewRestyClient(cfg)
	require.NotNil(t, cli)
	require.Equal(t, cfg.HostURL, cli.BaseURL)
	require.Equal(t, cfg.Timeout, cli.GetClient().Timeout)
}
