package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	dataRegex = "<script type=\"text\\/javascript\">window\\._sharedData = (.*);<\\/script>"
)

var r = regexp.MustCompile(dataRegex)

type Fetcher struct {
	hc       http.Client
	username string

	userAgent string
}

func NewFetcher(hc http.Client, username string) Fetcher {
	return Fetcher{
		hc:       hc,
		username: username,
	}
}

func (f Fetcher) fetchUserData() (Data, error) {
	html, err := f.fetchUserPage(f.username)
	if err != nil {
		return Data{}, err
	}

	jsonData, err := f.extractJsonData(html)

	if err != nil {
		return Data{}, err
	}

	var d Data
	err = json.Unmarshal(jsonData, &d)

	return d, err
}

func (f Fetcher) extractJsonData(html []byte) ([]byte, error) {
	rawData := r.FindSubmatch(html)

	if len(rawData) != 2 {
		return nil, errors.New("could not extract data from response")
	}

	return rawData[1], nil
}

func (f Fetcher) fetchUserPage(user string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseUrl, user), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := f.hc.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
