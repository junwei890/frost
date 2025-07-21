package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func normalizeURL(rawURL string) (string, error) {
	urlStructure, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return urlStructure.Host + strings.TrimRight(urlStructure.Path, "/"), nil
}

func urlsFromHTML(htmlBody, baseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	parseTree, err := html.Parse(reader)
	if err != nil {
		return []string{}, err
	}

	urls := []string{}
	for node := range parseTree.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, attribute := range node.Attr {
				if attribute.Key == "href" {
					urls = append(urls, attribute.Val)
				}
			}
		}
	}

	for index, rawURL := range urls {
		if urlStructure, err := url.Parse(rawURL); err != nil {
			return []string{}, err
		} else if urlStructure.Host == "" {
			urls[index] = baseURL + urlStructure.String()
		}
	}

	return urls, nil
}
