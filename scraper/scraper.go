package scraper

import (
	"log"
	"net/url"
	"sync"
)

type config struct {
	pages   map[string]int
	domain  *url.URL
	mu      *sync.Mutex
	control chan struct{}
	wg      *sync.WaitGroup
}

func InitiateCrawl(baseURL string) {
	domain, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	local := config{
		pages:   make(map[string]int),
		domain:  domain,
		mu:      &sync.Mutex{},
		control: make(chan struct{}, 5),
		wg:      &sync.WaitGroup{},
	}

	local.wg.Add(1)
	go local.crawlPage(baseURL)
	local.wg.Wait() // blocks till wait group is empty

	for key, value := range local.pages {
		log.Printf("%s: %d", key, value)
	}

}

func (c *config) visitCheck(normCurrURL string) bool {
	c.mu.Lock() // lock when reading or writing to map
	defer c.mu.Unlock()

	if _, ok := c.pages[normCurrURL]; ok {
		c.pages[normCurrURL] += 1
		return true
	}
	c.pages[normCurrURL] = 1
	return false
}

func (c *config) crawlPage(rawCurrURL string) {
	c.control <- struct{}{} // buffered channel limits number of requests/routines made
	defer func() {
		<-c.control
		c.wg.Done()
	}()

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
	if c.visitCheck(normCurrURL) { // checking if we already visited this site, return if yes
		return
	}

	html, err := getHTML(rawCurrURL)
	if err != nil {
		return
	}
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
