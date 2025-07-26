package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawUrl string) (string, error) {
	client := &http.Client{}
	res, err := client.Get(rawUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return "", errors.New("dead link")
	} else if 400 <= res.StatusCode && res.StatusCode < 500 {
		return "", errors.New("client error")
	} else if header := res.Header.Get("Content-Type"); !strings.Contains(header, "text/html") {
		return "", errors.New("content type not html")
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resData), nil
}

func normalizeURL(rawUrl string) (string, error) {
	urlStruct, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	return urlStruct.Host + strings.TrimRight(urlStruct.Path, "/"), nil
}
