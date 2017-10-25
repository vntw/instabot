package cli

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/venyii/instabot/provider"
	"github.com/venyii/instabot/slack"
)

type testSender struct {
	msgs []slack.Message
}

func (s *testSender) Send(msgs []slack.Message) {
	s.msgs = msgs
}

type testProvider struct {
	provider.Provider
}

type testDog struct{}

func (testDog) Sniff(p provider.Provider) ([]slack.Message, error) {
	return []slack.Message{
		{},
		{},
	}, nil
}

func (testDog) Sleep() {}

func TestDryRun(t *testing.T) {
	wd := testDog{}
	ig := testProvider{}
	sc := &testSender{}

	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Run(wd, ig, sc, true)

	actualOutput := buf.String()
	if !strings.Contains(actualOutput, "Looking for new media...") {
		t.Fatalf("unexpected log output %v", actualOutput)
	}
	if !strings.Contains(actualOutput, "dry-run enabled, would send 2 new messages") {
		t.Fatalf("unexpected log output %v", actualOutput)
	}
}
