package cfg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "TestNewConfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`
{
  "username": "user",
  "wait-time": 6,
  "slack-token": "slack__token-21",
  "slack-channel": "slack__channel-3",
  "proxy": ""
}
`)

	cfg, err := NewConfig(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if cfg.WaitTime != 6 {
		t.Fatalf("unexpected wait-time param, want 6 got %v", cfg.WaitTime)
	}
	if cfg.Username != "user" {
		t.Fatalf("unexpected username param, want user got %v", cfg.Username)
	}
	if cfg.SlackChannel != "slack__channel-3" {
		t.Fatalf("unexpected slack-channel param, want slack__channel-3 got %v", cfg.SlackChannel)
	}
	if cfg.SlackToken != "slack__token-21" {
		t.Fatalf("unexpected slack-token param, want slack__token-21 got %v", cfg.SlackToken)
	}
	if cfg.Proxy != "" {
		t.Fatalf("unexpected proxy param, want none got %v", cfg.Proxy)
	}
	if cfg.ProxyUrl() != nil {
		t.Fatalf("unexpected proxy url param, want nil got %v", cfg.ProxyUrl())
	}
}

func TestNewConfigParsesProxyUrl(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "TestNewConfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`
{
  "username": "u",
  "wait-time": 1,
  "slack-token": "st",
  "slack-channel": "sc",
  "proxy": "http://127.0.0.1:8080"
}
`)

	cfg, err := NewConfig(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Proxy != "http://127.0.0.1:8080" {
		t.Fatalf("unexpected proxy param, want http://127.0.0.1:8080 got %v", cfg.Proxy)
	}
	if cfg.ProxyUrl() == nil {
		t.Fatal("unexpected nil proxy url")
	}
}

func resetFile(f *os.File) {
	f.Truncate(0)
	f.Seek(0, 0)
}

func TestNewConfigValidatesParams(t *testing.T) {
	var err error

	tmpFile, err := ioutil.TempFile("", "TestNewConfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// username
	tmpFile.WriteString(`
{
  "username": "",
  "wait-time": 1,
  "slack-token": "st",
  "slack-channel": "sc",
  "proxy": ""
}
`)

	_, err = NewConfig(tmpFile.Name())
	if err == nil {
		t.Fatal("expected config error, got nil")
	}

	if err.Error() != "invalid config: Username must not be empty" {
		t.Fatalf("unexpected config error, got %v", err)
	}

	// wait time
	resetFile(tmpFile)
	tmpFile.WriteString(`
{
  "username": "u",
  "wait-time": -1,
  "slack-token": "st",
  "slack-channel": "sc",
  "proxy": ""
}
`)

	_, err = NewConfig(tmpFile.Name())
	if err == nil {
		t.Fatal("expected config error, got nil")
	}

	if err.Error() != "could not unmarshal config: json: cannot unmarshal number -1 into Go struct field Config.wait-time of type uint8" {
		t.Fatalf("unexpected config error, got %v", err)
	}

	// slack token
	resetFile(tmpFile)
	tmpFile.WriteString(`
{
  "username": "u",
  "wait-time": 1,
  "slack-token": "",
  "slack-channel": "sc",
  "proxy": ""
}
`)

	_, err = NewConfig(tmpFile.Name())
	if err == nil {
		t.Fatal("expected config error, got nil")
	}

	if err.Error() != "invalid config: Slack not configured" {
		t.Fatalf("unexpected config error, got %v", err)
	}

	// slack channel
	resetFile(tmpFile)
	tmpFile.WriteString(`
{
  "username": "u",
  "wait-time": 1,
  "slack-token": "st",
  "slack-channel": "",
  "proxy": ""
}
`)

	_, err = NewConfig(tmpFile.Name())
	if err == nil {
		t.Fatal("expected config error, got nil")
	}

	if err.Error() != "invalid config: Slack not configured" {
		t.Fatalf("unexpected config error, got %v", err)
	}
}
