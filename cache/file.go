package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func NewFileCache(ns string) *fileCache {
	return &fileCache{filename: fmt.Sprintf("%s-cache.txt", ns)}
}

type fileCache struct {
	filename string
}

func (c fileCache) ReadLastDate() int64 {
	val, err := ioutil.ReadFile(c.filename)

	if err != nil {
		if os.IsNotExist(err) {
			return -1
		} else {
			log.Println("Could not read from cache")
			log.Fatal(err)
		}
	}

	str := strings.TrimSpace(string(val))
	lastDate, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		// Invalid date value given - ignore
		return -1
	}

	return lastDate
}

func (c fileCache) WriteLastDate(date int64) error {
	return ioutil.WriteFile(c.filename, []byte(strconv.FormatInt(date, 10)), 0644)
}
