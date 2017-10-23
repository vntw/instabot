package instagram

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/venyii/instabot/provider/cache"
)

type dummyCache struct {
	cache.Cache
}

type testTransport struct{}

func (t testTransport) RoundTrip(*http.Request) (*http.Response, error) {
	f, err := os.Open("fixtures/response.html")
	if err != nil {
		log.Fatal(err)
	}

	r := &http.Response{
		Body: ioutil.NopCloser(f),
	}

	return r, nil
}

func TestLatest(t *testing.T) {
	a := http.Client{Transport: testTransport{}}

	p := NewProvider(dummyCache{}, a, "test", "#ff000")
	msgs, err := p.Latest()

	if err != nil {
		log.Fatal(err)
	}

	if len(msgs) != 12 {
		t.Fatalf("unexpected msgs count, want 12 got %v", len(msgs))
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
		1505237405,
		1504652223,
		1503969273,
		1503698278,
	}

	for i := 0; i < len(expectedIds); i++ {
		if msgs[i].Date != expectedIds[i] {
			t.Fatalf("unexpected message returned at index %d", i)
		}
	}
}
