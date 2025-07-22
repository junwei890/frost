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

func getHTML(rawURL string) (string, error) { // get request for full page html
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

func normalizeURL(rawURL string) (string, error) { // helper function for checking page visits
	urlStruct, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return urlStruct.Host + strings.TrimRight(urlStruct.Path, "/"), nil
}

func urlsFromHTML(htmlBody string, host *url.URL) ([]string, error) {
	htmlTree, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	urls := []string{} // extracting all links from <a> tags
	for node := range htmlTree.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					urls = append(urls, attr.Val)
				}
			}
		}
	}

	for i, rawURL := range urls { // and cleaning up if there's no host name
		if urlStruct, err := url.Parse(rawURL); err != nil {
			return []string{}, err
		} else if urlStruct.Hostname() == "" {
			urls[i] = host.ResolveReference(urlStruct).String()
		}
	}

	return urls, nil
}
