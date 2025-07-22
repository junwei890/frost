package scraper

import (
	"log"
	"net/url"
)

func CrawlPage(domain *url.URL, rawCurrURL string, pages map[string]int) {
	currStruct, err := url.Parse(rawCurrURL)
	if err != nil {
		return
	}
	if domain.Hostname() != currStruct.Hostname() {
		return
	}

	normCurrURL, err := normalizeURL(rawCurrURL)
	if err != nil {
		return
	}
	if _, ok := pages[normCurrURL]; ok {
		pages[normCurrURL]++
		return
	}
	pages[normCurrURL] = 1

	html, err := getHTML(rawCurrURL)
	if err != nil {
		return
	}
	links, err := urlsFromHTML(html, rawCurrURL)
	if err != nil {
		return
	}
	for _, link := range links {
		log.Printf("crawling %s", link)
		CrawlPage(domain, link, pages)
	}
}
