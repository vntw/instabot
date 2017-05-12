package cache

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const saveFile = "lastdate.txt"

func NewFileCache() Cache {
	return &FileCache{}
}

type FileCache struct {
}

func (c FileCache) ReadLastDate() int64 {
	val, err := ioutil.ReadFile(saveFile)

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

func (c FileCache) WriteLastDate(date int64) error {
	return ioutil.WriteFile(saveFile, []byte(strconv.FormatInt(date, 10)), 0644)
}
