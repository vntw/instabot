package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func NewFileCache(ns string, u string) *fileCache {
	return &fileCache{filename: fmt.Sprintf("%s-%s-cache.txt", ns, u)}
}

type fileCache struct {
	filename string
}

func (c fileCache) ReadLastDate() (int64, error) {
	val, err := ioutil.ReadFile(c.filename)

	if err != nil {
		if os.IsNotExist(err) {
			return -1, nil
		}

		return 0, errors.New(fmt.Sprintf("could not read from cache file %s: %v", c.filename, err))
	}

	str := strings.TrimSpace(string(val))
	lastDate, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		// Invalid date value given - ignore
		return -1, nil
	}

	return lastDate, nil
}

func (c fileCache) WriteLastDate(date int64) error {
	return ioutil.WriteFile(c.filename, []byte(strconv.FormatInt(date, 10)), 0644)
}
