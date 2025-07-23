package scraper

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	client := &http.Client{}
	res, err := client.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return "", errors.New("dead link")
	} else if 400 <= res.StatusCode && res.StatusCode < 500 {
		return "", errors.New("client error")
	}
	if header := res.Header.Get("Content-Type"); !strings.Contains(header, "text/html") {
		return "", errors.New("content type not html")
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resData), nil
}

func normalizeURL(rawURL string) (string, error) {
	urlStruct, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return urlStruct.Host + strings.TrimRight(urlStruct.Path, "/"), nil
}
