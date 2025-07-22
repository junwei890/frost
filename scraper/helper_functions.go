package scraper

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func urlsFromHTML(htmlBody, host string) ([]string, error) {
	htmlTree, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	urls := []string{}
	for node := range htmlTree.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					urls = append(urls, attr.Val)
				}
			}
		}
	}

	for i, rawURL := range urls {
		if urlStruct, err := url.Parse(rawURL); err != nil {
			return []string{}, err
		} else if urlStruct.Hostname() == "" {
			urls[i] = host + urlStruct.String()
		}
	}

	return urls, nil
}
