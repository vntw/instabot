package instagram

import (
	"log"
	"testing"

	"github.com/venyii/instabot/cfg"
)

type testCache struct {
	Date int64
}

func (c *testCache) ReadLastDate() int64 {
	return c.Date
}

func (c *testCache) WriteLastDate(date int64) error {
	c.Date = date
	return nil
}

func TestLatest(t *testing.T) {
	cache := &testCache{Date: 1505237479}
	config := cfg.Config{}
	p := NewProvider(config, cache)
	msgs, err := p.Latest(false, true)

	if err != nil {
		log.Fatal(err)
	}

	if len(msgs) != 8 {
		t.Fatalf("unexpected msgs count, want 8 got %v", len(msgs))
	}

	expectedIds := []int64{
		1507829725,
		1507747423,
		1507575802,
		1507355136,
		1506841794,
		1506381081,
		1505931970,
		1505413351,
	}

	for i := 0; i < 8; i++ {
		if msgs[i].Date != expectedIds[i] {
			t.Fatalf("unexpected message returned at index %d", i)
		}
	}

	if cache.Date != 1507829725 {
		t.Fatalf("unexpected cache date, want 1507829725 got %v", cache.Date)
	}
}
