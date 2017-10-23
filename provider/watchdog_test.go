package provider

import (
	"testing"

	"github.com/venyii/instabot/provider/cache"
	"github.com/venyii/instabot/slack"
)

type testCache struct {
	Date int64
}

func (c *testCache) ReadLastDate() (int64, error) {
	return c.Date, nil
}

func (c *testCache) WriteLastDate(date int64) error {
	c.Date = date
	return nil
}

type dummyProvider struct {
	cache cache.Cache
}

func (dummyProvider) Latest() ([]slack.Message, error) {
	return []slack.Message{
		{Date: 110},
		{Date: 108},
		{Date: 103},
	}, nil
}

func (p dummyProvider) Cache() cache.Cache {
	return p.cache
}

func TestSniffWithOldCache(t *testing.T) {
	c := &testCache{101}
	p := dummyProvider{cache: c}

	wd := NewWatchDog(1)
	msgs, err := wd.Sniff(p)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(msgs) != 3 {
		t.Fatalf("unexpected message count, want 3 got %v", len(msgs))
	}

	if c.Date != 110 {
		t.Fatalf("unexpected cache date, want 110 got %v", c.Date)
	}
}

func TestSniffWithSomewhatOldCache(t *testing.T) {
	c := &testCache{105}
	p := dummyProvider{cache: c}

	wd := NewWatchDog(1)
	msgs, err := wd.Sniff(p)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("unexpected message count, want 2 got %v", len(msgs))
	}

	if c.Date != 110 {
		t.Fatalf("unexpected cache date, want 110 got %v", c.Date)
	}
}

func TestSniffWithNewCache(t *testing.T) {
	c := &testCache{110}
	p := dummyProvider{cache: c}

	wd := NewWatchDog(1)
	msgs, err := wd.Sniff(p)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(msgs) != 0 {
		t.Fatalf("unexpected message count, want 0 got %v (%v)", len(msgs), msgs)
	}

	if c.Date != 110 {
		t.Fatalf("unexpected cache date, want 110 got %v", c.Date)
	}
}
