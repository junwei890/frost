package scraper

import (
	"log"
	"net/url"
	"sync"
)

type data struct {
	title string
	count int
}

type config struct {
	metadata  map[string]data
	domain    *url.URL
	mu        *sync.Mutex
	control   chan struct{}
	wg        *sync.WaitGroup
	maxVisits int
}

func InitiateCrawl(baseURL string) {
	domain, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	local := config{
		metadata:  make(map[string]data),
		domain:    domain,
		mu:        &sync.Mutex{},
		control:   make(chan struct{}, 5),
		wg:        &sync.WaitGroup{},
		maxVisits: 20,
	}

	local.wg.Add(1)
	go local.crawlPage(baseURL)
	local.wg.Wait() // blocks till wait group is empty

	for key, value := range local.metadata {
		log.Printf("%s: %s", key, value.title)
	}
}

func (c *config) urlVisited(normCurrURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.metadata[normCurrURL]; ok {
		c.metadata[normCurrURL] = data{
			title: c.metadata[normCurrURL].title,
			count: c.metadata[normCurrURL].count + 1,
		}
		return true
	}
	c.metadata[normCurrURL] = data{
		count: 1,
	}
	return false
}

func (c *config) setTitle(normCurrURL, title string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metadata[normCurrURL] = data{
		title: title,
		count: c.metadata[normCurrURL].count,
	}
}

func (c *config) maxReached() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.metadata) >= c.maxVisits {
		return true
	}
	return false
}

func (c *config) crawlPage(rawCurrURL string) {
	c.control <- struct{}{} // buffered channel limits number of requests/routines made
	defer func() {
		<-c.control
		c.wg.Done()
	}()

	if c.maxReached() {
		return
	}

	currStruct, err := url.Parse(rawCurrURL)
	if err != nil {
		return
	}
	if c.domain.Hostname() != currStruct.Hostname() { // only want to scrape within given domain
		return
	}

	normCurrURL, err := normalizeURL(rawCurrURL)
	if err != nil {
		return
	}
	if c.urlVisited(normCurrURL) { // checking if we already visited this site, return if yes
		return
	}

	html, err := getHTML(rawCurrURL)
	if err != nil {
		return
	}

	title, err := titleFromHTML(html)
	if err != nil {
		return
	}
	c.setTitle(normCurrURL, title)

	links, err := urlsFromHTML(html, c.domain) // passing in domain name in case hrefs are missing hostname
	if err != nil {
		return
	}
	for _, link := range links {
		c.wg.Add(1)
		log.Printf("crawling %s", link)
		go c.crawlPage(link)
	}
}
