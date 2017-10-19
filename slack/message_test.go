package slack

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage("custom", "user")
	if msg.ProviderName != "custom" {
		t.Fatalf("unexpected provider name, want custom got %v", msg.ProviderName)
	}
}

func TestMessageSlackMessageParams(t *testing.T) {
	msg := NewMessage("custom", "user")
	msg.Author = "author"
	msg.AuthorIcon = "author_icon"
	msg.AuthorLink = "/u/author"
	msg.Color = "#fddd4a"
	msg.Date = 1508271836
	msg.ImageUrl = "/i/image.jpg"
	msg.Text = "Hello World"
	msg.Fallback = "Hello World Fallback"

	params := msg.SlackMessageParams()
	if len(params.Attachments) != 1 {
		t.Fatalf("unexpected amount of attachments, want 1 got %v", len(params.Attachments))
	}
	a := params.Attachments[0]
	if a.Ts.String() != "1508271836" {
		t.Fatalf("unexpected timestamp, want 1508271836 got %v", a.Ts.String())
	}
}
