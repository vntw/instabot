package instagram

import (
	"testing"
)

func TestUserMediaNodes(t *testing.T) {
	u := User{
		Username: "u",
		Media: Media{
			Nodes: []Node{
				{Id: "id1", Src: "src1", Code: "code1", Date: 1337},
				{Id: "id2", Src: "src2", Code: "code2", Date: 1338},
			},
		},
	}
	n := u.MediaNodes(1)

	if len(n) != 1 {
		t.Fatalf("unexpected length of user media nodes, want 1 got %v", len(n))
	}
	if n[0].Id != "id1" {
		t.Fatalf("unexpected media node returned, want id1 got %v", n[0].Id)
	}
}

func TestUserProfileUrl(t *testing.T) {
	u := User{Username: "user"}

	if u.ProfileUrl() != "https://www.instagram.com/user" {
		t.Fatalf("unexpected user profile url, want https://www.instagram.com/user got %v", u.ProfileUrl())
	}
}

func TestNodeDetailUrl(t *testing.T) {
	n := Node{Id: "id1", Src: "src1", Code: "code1", Date: 1337}

	if n.DetailUrl() != "https://www.instagram.com/p/code1" {
		t.Fatalf("unexpected node detail url, want https://www.instagram.com/p/code1 got %v", n.DetailUrl())
	}
}
