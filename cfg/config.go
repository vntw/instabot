package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

type Config struct {
	Username string `json:"username"`
	// The time between requests in minutes
	WaitTime     uint8  `json:"wait-time"`
	SlackToken   string `json:"slack-token"`
	SlackChannel string `json:"slack-channel"`
	Proxy        string `json:"proxy"`

	proxyUrl *url.URL
}

func NewConfig(configFile string) (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("could not read config file: " + err.Error())
	}

	var cfg *Config
	if err = json.Unmarshal(content, &cfg); err != nil {
		return nil, errors.New("could not unmarshal config: " + err.Error())
	}

	if cfg.Username == "" {
		return nil, errors.New("invalid config: Username must not be empty")
	}
	if cfg.WaitTime <= 0 {
		return nil, errors.New("invalid config: WaitTime must be greater than zero")
	}
	if cfg.SlackToken == "" || cfg.SlackChannel == "" {
		return nil, errors.New("invalid config: Slack not configured")
	}

	if err := cfg.populateProxyUrl(); err != nil {
		return nil, errors.New("could not create proxy url: " + err.Error())
	}

	return cfg, nil
}

func (c *Config) ProxyUrl() *url.URL {
	return c.proxyUrl
}

func (c *Config) populateProxyUrl() error {
	if c.Proxy == "" {
		return nil
	}

	proxyUrl, err := url.Parse(c.Proxy)
	if err != nil {
		return err
	}

	c.proxyUrl = proxyUrl

	return nil
}
