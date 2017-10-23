package cache

import (
	"io/ioutil"
	"os"
	"testing"
)

const testCacheFilename = "namespace-username-cache.txt"

func TestReadLastDateWithNoCacheFile(t *testing.T) {
	c := NewFileCache("namespace", "username")

	ld, err := c.ReadLastDate()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if ld != -1 {
		t.Fatalf("unexpected lastDate, want -1 got %v", ld)
	}
}

func TestReadLastDateWithCacheFile(t *testing.T) {
	c := NewFileCache("namespace", "username")

	if err := ioutil.WriteFile(testCacheFilename, []byte("110"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testCacheFilename)

	ld, err := c.ReadLastDate()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if ld != 110 {
		t.Fatalf("unexpected lastDate, want 110 got %v", ld)
	}
}

func TestWriteLastDate(t *testing.T) {
	c := NewFileCache("namespace", "username")

	err := c.WriteLastDate(110)
	if err != nil {
		t.Fatalf("unexpected error writing cache file: %v", err)
	}
	defer os.Remove(testCacheFilename)

	ld, err := ioutil.ReadFile(testCacheFilename)
	if err != nil {
		t.Fatal(err)
	}

	if string(ld) != "110" {
		t.Fatalf("unexpected lastDate in cache file, want 110 got %v", string(ld))
	}
}
