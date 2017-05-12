package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	dataRegex = "<script type=\"text\\/javascript\">window\\._sharedData = (.*);<\\/script>"
)

var r = regexp.MustCompile(dataRegex)

type Fetcher struct {
	proxyUrl *url.URL

	userAgent string
}

func NewFetcher(p *url.URL) Fetcher {
	return Fetcher{proxyUrl: p}
}

func (f Fetcher) FetchUserData(user string) (Data, error) {
	html, err := f.fetchUserPage(user)
	if err != nil {
		return Data{}, err
	}

	return f.ExtractDataFromHtml(html)
}

func (f Fetcher) ExtractDataFromHtml(html []byte) (Data, error) {
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
		log.Printf("Unexpected raw data: %#v\n", rawData)
		return nil, errors.New("could not extract data from response")
	}

	return rawData[1], nil
}

func (f Fetcher) fetchUserPage(user string) ([]byte, error) {
	var client *http.Client

	if f.proxyUrl != nil {
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(f.proxyUrl)}}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseUrl, user), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
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
