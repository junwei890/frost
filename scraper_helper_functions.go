package main

import (
	"errors"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	urlStructure, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("can't parse invalid url")
	}
	return urlStructure.Host + strings.TrimRight(urlStructure.Path, "/"), nil
}
